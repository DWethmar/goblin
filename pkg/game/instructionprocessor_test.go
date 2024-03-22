package game

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain/actor"
	"github.com/dwethmar/goblin/pkg/services"
)

func TestInstructionProcessor_use(t *testing.T) {
	type fields struct {
		Logger       *slog.Logger
		AggregateID  string
		ActorService ActorService
	}
	type args struct {
		instruction string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				Logger:       slog.Default(),
				AggregateID:  "1",
				ActorService: nil,
			},
			args: args{
				instruction: "use 1",
			},
			wantErr: false,
		},
		{
			name: "should return error if instruction is invalid",
			fields: fields{
				Logger:       slog.Default(),
				AggregateID:  "1",
				ActorService: nil,
			},
			args: args{
				instruction: "use",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstructionProcessor{
				Logger:       tt.fields.Logger,
				AggregateID:  tt.fields.AggregateID,
				ActorService: tt.fields.ActorService,
			}
			if err := p.use(tt.args.instruction); (err != nil) != tt.wantErr {
				t.Errorf("InstructionProcessor.use() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstructionProcessor_create(t *testing.T) {
	type fields struct {
		Logger       *slog.Logger
		AggregateID  string
		ActorService ActorService
	}
	type args struct {
		ctx context.Context
		in  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				Logger:      slog.Default(),
				AggregateID: "1",
				ActorService: services.NewActorService(
					&actor.MockRepository{
						GetFunc: func(ctx context.Context, id string) (*actor.Actor, error) {
							return nil, errors.New("could not find actor")
						},
					},
					aggr.NoopCommandBus,
				),
			},
			args: args{
				ctx: context.Background(),
				in:  "create actor test 1 1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstructionProcessor{
				Logger:       tt.fields.Logger,
				AggregateID:  tt.fields.AggregateID,
				ActorService: tt.fields.ActorService,
			}
			if err := p.create(tt.args.ctx, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("InstructionProcessor.create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstructionProcessor_createActor(t *testing.T) {
	type fields struct {
		Logger       *slog.Logger
		AggregateID  string
		ActorService ActorService
	}
	type args struct {
		ctx context.Context
		in  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstructionProcessor{
				Logger:       tt.fields.Logger,
				AggregateID:  tt.fields.AggregateID,
				ActorService: tt.fields.ActorService,
			}
			if err := p.createActor(tt.args.ctx, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("InstructionProcessor.createActor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstructionProcessor_move(t *testing.T) {
	type fields struct {
		Logger       *slog.Logger
		AggregateID  string
		ActorService ActorService
	}
	type args struct {
		ctx context.Context
		in  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstructionProcessor{
				Logger:       tt.fields.Logger,
				AggregateID:  tt.fields.AggregateID,
				ActorService: tt.fields.ActorService,
			}
			if err := p.move(tt.args.ctx, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("InstructionProcessor.move() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstructionProcessor_Process(t *testing.T) {
	type fields struct {
		Logger       *slog.Logger
		AggregateID  string
		ActorService ActorService
	}
	type args struct {
		ctx context.Context
		in  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstructionProcessor{
				Logger:       tt.fields.Logger,
				AggregateID:  tt.fields.AggregateID,
				ActorService: tt.fields.ActorService,
			}
			if err := p.Process(tt.args.ctx, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("InstructionProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
