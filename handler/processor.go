package handler

import (
	"context"
	metaerror "meta/meta-error"
)

type Processor[T any] struct {
	handlers []Handler[T]
}

func (p *Processor[T]) Init() error {
	for _, handler := range p.handlers {
		err := handler.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor[T]) Process(ctx context.Context, event T) error {
	var finalErr error
	for _, handler := range p.handlers {
		if handler.IsShouldProcess(ctx, event) {
			continueProcess, err := handler.DoProcess(ctx, event)
			if err != nil {
				finalErr = metaerror.Join(finalErr, err)
			}
			if !continueProcess {
				break
			}
		}
	}
	return finalErr
}

type ProcessorBuilder[T any] struct {
	processor *Processor[T]
}

func NewProcessorBuilder[T any]() *ProcessorBuilder[T] {
	return &ProcessorBuilder[T]{processor: &Processor[T]{}}
}

func (builder *ProcessorBuilder[T]) Build() (*Processor[T], error) {
	err := builder.processor.Init()
	if err != nil {
		return nil, err
	}
	return builder.processor, nil
}

func (builder *ProcessorBuilder[T]) AddHandler(handler Handler[T]) *ProcessorBuilder[T] {
	if handler != nil {
		builder.processor.handlers = append(builder.processor.handlers, handler)
	}
	return builder
}
