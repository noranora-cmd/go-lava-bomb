package volcano

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_checkFile(t *testing.T) {
	type args struct {
		confCheckTime int64
		configFile    string
	}
	tests := []struct {
		name    string
		args    args
		want    *noise
		wantErr bool
	}{
		{
			name: "File exists with all values",
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test.json",
			},
			want: &noise{
				Second: "rumble",
				Minute: "rrrumble",
				Hour:   "boom",
			},
			wantErr: false,
		},
		{
			name: "File exists with some values missing",
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test1.json",
			},
			want: &noise{
				Second: "rumble",
				Minute: "rrrumble",
				Hour:   "",
			},
			wantErr: false,
		},
		{
			name: "File exists with incorrect json",
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test2.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "File exists with some values incorrect",
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test3.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "File does not exist",
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test100.json",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkFile(tt.args.confCheckTime, tt.args.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

type monitorOutput struct {
	secNoise, minNoise, hourNoise string
}

func Test_volcano_monitor(t *testing.T) {
	type fields struct {
		secCh         chan string
		minCh         chan string
		hourCh        chan string
		confCheckTime int64
		errorCh       chan error
	}
	type args struct {
		confCheckTime int64
		configFile    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    monitorOutput
		wantErr bool
	}{
		{
			name: "missing config file",
			fields: fields{
				secCh:         make(chan string, 1),
				minCh:         make(chan string, 1),
				hourCh:        make(chan string, 1),
				confCheckTime: 1,
				errorCh:       make(chan error, 1),
			},
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test100.json",
			},
			want:    monitorOutput{},
			wantErr: true,
		},
		{
			name: "config file with incorrect json format",
			fields: fields{
				secCh:         make(chan string, 1),
				minCh:         make(chan string, 1),
				hourCh:        make(chan string, 1),
				confCheckTime: 1,
				errorCh:       make(chan error, 1),
			},
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test2.json",
			},
			want:    monitorOutput{},
			wantErr: true,
		},
		{
			name: "config file with incorrect noise values",
			fields: fields{
				secCh:         make(chan string, 1),
				minCh:         make(chan string, 1),
				hourCh:        make(chan string, 1),
				confCheckTime: 1,
				errorCh:       make(chan error, 1),
			},
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test3.json",
			},
			want:    monitorOutput{},
			wantErr: true,
		},
		{
			name: "correct config file with all values",
			fields: fields{
				secCh:         make(chan string, 1),
				minCh:         make(chan string, 1),
				hourCh:        make(chan string, 1),
				confCheckTime: 1,
				errorCh:       make(chan error, 1),
			},
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test.json",
			},
			want: monitorOutput{
				secNoise:  "rumble",
				minNoise:  "rrrumble",
				hourNoise: "boom",
			},
			wantErr: false,
		},
		{
			name: "correct config file with some values",
			fields: fields{
				secCh:         make(chan string, 1),
				minCh:         make(chan string, 1),
				hourCh:        make(chan string, 1),
				confCheckTime: 1,
				errorCh:       make(chan error, 1),
			},
			args: args{
				confCheckTime: 1,
				configFile:    "testdata/config_test1.json",
			},
			want: monitorOutput{
				secNoise:  "rumble",
				minNoise:  "rrrumble",
				hourNoise: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &volcano{
				secCh:         tt.fields.secCh,
				minCh:         tt.fields.minCh,
				hourCh:        tt.fields.hourCh,
				confCheckTime: tt.fields.confCheckTime,
				errorCh:       tt.fields.errorCh,
			}
			err := c.monitor(tt.args.confCheckTime, tt.args.configFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("checkFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			var secNoise, minNoise, hourNoise string
			if len(c.secCh) > 0 {
				secNoise = <-c.secCh
			}
			if len(c.minCh) > 0 {
				minNoise = <-c.minCh
			}
			if len(c.hourCh) > 0 {
				hourNoise = <-c.hourCh
			}

			got := monitorOutput{
				secNoise:  secNoise,
				minNoise:  minNoise,
				hourNoise: hourNoise,
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

type rumbleOutput struct {
	noises []string
}

func Test_volcano_rumble(t *testing.T) {
	type fields struct {
		secCh         chan string
		minCh         chan string
		hourCh        chan string
		stopTime      time.Duration
		confCheckTime int64
		doneCh        chan struct{}
		errorCh       chan error
		noiseCh       chan string
		mu            sync.WaitGroup

		secNoisesToSend []string
		minNoisesToSend []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    rumbleOutput
		wantErr bool
	}{
		{
			name: "10 second rumbleing with no dynamic noise changes",
			fields: fields{
				secCh:         make(chan string, 1),
				minCh:         make(chan string, 1),
				hourCh:        make(chan string, 1),
				stopTime:      10 * time.Second,
				confCheckTime: 1,
				doneCh:        make(chan struct{}, 1),
				errorCh:       make(chan error, 1),
				noiseCh:       make(chan string),

				secNoisesToSend: []string{"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble"},
				minNoisesToSend: []string{},
			},
			want: rumbleOutput{
				noises: []string{"1  ...  rumble", "2  ...  rumble", "3  ...  rumble", "4  ...  rumble", "5  ...  rumble", "6  ...  rumble", "7  ...  rumble", "8  ...  rumble", "9  ...  rumble", "10  ...  rumble"},
			},
			wantErr: false,
		},
		{
			name: "61 second rumbling with no dynamic noise changes",
			fields: fields{
				secCh:         make(chan string, 1),
				minCh:         make(chan string, 1),
				hourCh:        make(chan string, 1),
				stopTime:      61 * time.Second,
				confCheckTime: 1,
				doneCh:        make(chan struct{}, 1),
				errorCh:       make(chan error, 1),
				noiseCh:       make(chan string),

				secNoisesToSend: []string{
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
				},
				minNoisesToSend: []string{},
			},
			want: rumbleOutput{
				noises: []string{
					"1  ...  rumble", "2  ...  rumble", "3  ...  rumble", "4  ...  rumble", "5  ...  rumble", "6  ...  rumble", "7  ...  rumble", "8  ...  rumble", "9  ...  rumble", "10  ...  rumble",
					"11  ...  rumble", "12  ...  rumble", "13  ...  rumble", "14  ...  rumble", "15  ...  rumble", "16  ...  rumble", "17  ...  rumble", "18  ...  rumble", "19  ...  rumble", "20  ...  rumble",
					"21  ...  rumble", "22  ...  rumble", "23  ...  rumble", "24  ...  rumble", "25  ...  rumble", "26  ...  rumble", "27  ...  rumble", "28  ...  rumble", "29  ...  rumble", "30  ...  rumble",
					"31  ...  rumble", "32  ...  rumble", "33  ...  rumble", "34  ...  rumble", "35  ...  rumble", "36  ...  rumble", "37  ...  rumble", "38  ...  rumble", "39  ...  rumble", "40  ...  rumble",
					"41  ...  rumble", "42  ...  rumble", "43  ...  rumble", "44  ...  rumble", "45  ...  rumble", "46  ...  rumble", "47  ...  rumble", "48  ...  rumble", "49  ...  rumble", "50  ...  rumble",
					"51  ...  rumble", "52  ...  rumble", "53  ...  rumble", "54  ...  rumble", "55  ...  rumble", "56  ...  rumble", "57  ...  rumble", "58  ...  rumble", "59  ...  rumble", "60  ...  RUMBLE", "61  ...  rumble",
				},
			},
			wantErr: false,
		},
		{
			name: "61 second rumbling with dynamic noise changes",
			fields: fields{
				secCh:         make(chan string, 1),
				minCh:         make(chan string, 1),
				hourCh:        make(chan string, 1),
				stopTime:      61 * time.Second,
				confCheckTime: 1,
				doneCh:        make(chan struct{}, 1),
				errorCh:       make(chan error, 1),
				noiseCh:       make(chan string),

				secNoisesToSend: []string{
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
					"rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble", "rumble",
				},
				minNoisesToSend: []string{"rrrumble"},
			},
			want: rumbleOutput{
				noises: []string{
					"1  ...  rumble", "2  ...  rumble", "3  ...  rumble", "4  ...  rumble", "5  ...  rumble", "6  ...  rumble", "7  ...  rumble", "8  ...  rumble", "9  ...  rumble", "10  ...  rumble",
					"11  ...  rumble", "12  ...  rumble", "13  ...  rumble", "14  ...  rumble", "15  ...  rumble", "16  ...  rumble", "17  ...  rumble", "18  ...  rumble", "19  ...  rumble", "20  ...  rumble",
					"21  ...  rumble", "22  ...  rumble", "23  ...  rumble", "24  ...  rumble", "25  ...  rumble", "26  ...  rumble", "27  ...  rumble", "28  ...  rumble", "29  ...  rumble", "30  ...  rumble",
					"31  ...  rumble", "32  ...  rumble", "33  ...  rumble", "34  ...  rumble", "35  ...  rumble", "36  ...  rumble", "37  ...  rumble", "38  ...  rumble", "39  ...  rumble", "40  ...  rumble",
					"41  ...  rumble", "42  ...  rumble", "43  ...  rumble", "44  ...  rumble", "45  ...  rumble", "46  ...  rumble", "47  ...  rumble", "48  ...  rumble", "49  ...  rumble", "50  ...  rumble",
					"51  ...  rumble", "52  ...  rumble", "53  ...  rumble", "54  ...  rumble", "55  ...  rumble", "56  ...  rumble", "57  ...  rumble", "58  ...  rumble", "59  ...  rumble", "60  ...  rrrumble", "61  ...  rumble",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &volcano{
				secCh:         tt.fields.secCh,
				minCh:         tt.fields.minCh,
				hourCh:        tt.fields.hourCh,
				stopTime:      tt.fields.stopTime,
				confCheckTime: tt.fields.confCheckTime,
				doneCh:        tt.fields.doneCh,
				errorCh:       tt.fields.errorCh,
				noiseCh:       tt.fields.noiseCh,
				mu:            tt.fields.mu,
			}
			go func() {
				noiseMax := len(tt.fields.secNoisesToSend)
				for i := 0; i < noiseMax; i++ {
					c.secCh <- tt.fields.secNoisesToSend[i]
				}
			}()
			go func() {
				time.Sleep(5 * time.Second)
				noiseMax := len(tt.fields.minNoisesToSend)
				for i := 0; i < noiseMax; i++ {
					c.minCh <- tt.fields.minNoisesToSend[i]
				}
			}()

			c.mu.Add(1)
			go c.rumble()

			var err error
			if tt.wantErr {
				err = <-c.errorCh
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("rumble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			gotNoises := []string{}
			for msg := range c.noiseCh {
				gotNoises = append(gotNoises, msg)
			}

			if !reflect.DeepEqual(gotNoises, tt.want.noises) {
				t.Errorf("got = %v", gotNoises)
				t.Errorf("want %v", tt.want.noises)
			}
		})
	}
}
