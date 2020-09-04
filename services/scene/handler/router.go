// Code generated by jrpc. DO NOT EDIT.

package handler

import (
	context "context"

	taxi "github.com/jakewright/home-automation/libraries/go/taxi"
	def "github.com/jakewright/home-automation/services/scene/def"
)

// taxiRouter is an interface implemented by taxi.Router
type taxiRouter interface {
	HandleFunc(method, path string, handler func(context.Context, taxi.Decoder) (interface{}, error))
}

type handler interface {
	CreateScene(ctx context.Context, body *def.CreateSceneRequest) (*def.CreateSceneResponse, error)
	ReadScene(ctx context.Context, body *def.ReadSceneRequest) (*def.ReadSceneResponse, error)
	ListScenes(ctx context.Context, body *def.ListScenesRequest) (*def.ListScenesResponse, error)
	DeleteScene(ctx context.Context, body *def.DeleteSceneRequest) (*def.DeleteSceneResponse, error)
	SetScene(ctx context.Context, body *def.SetSceneRequest) (*def.SetSceneResponse, error)
}

// RegisterRoutes adds the service's routes to the router
func RegisterRoutes(r taxiRouter, h handler) {
	r.HandleFunc("POST", "/scenes", func(ctx context.Context, decode taxi.Decoder) (interface{}, error) {
		body := &def.CreateSceneRequest{}
		if err := decode(body); err != nil {
			return nil, err
		}

		if err := body.Validate(); err != nil {
			return nil, err
		}

		return h.CreateScene(ctx, body)
	})

	r.HandleFunc("GET", "/scene", func(ctx context.Context, decode taxi.Decoder) (interface{}, error) {
		body := &def.ReadSceneRequest{}
		if err := decode(body); err != nil {
			return nil, err
		}

		if err := body.Validate(); err != nil {
			return nil, err
		}

		return h.ReadScene(ctx, body)
	})

	r.HandleFunc("GET", "/scenes", func(ctx context.Context, decode taxi.Decoder) (interface{}, error) {
		body := &def.ListScenesRequest{}
		if err := decode(body); err != nil {
			return nil, err
		}

		if err := body.Validate(); err != nil {
			return nil, err
		}

		return h.ListScenes(ctx, body)
	})

	r.HandleFunc("DELETE", "/scene", func(ctx context.Context, decode taxi.Decoder) (interface{}, error) {
		body := &def.DeleteSceneRequest{}
		if err := decode(body); err != nil {
			return nil, err
		}

		if err := body.Validate(); err != nil {
			return nil, err
		}

		return h.DeleteScene(ctx, body)
	})

	r.HandleFunc("POST", "/scene/set", func(ctx context.Context, decode taxi.Decoder) (interface{}, error) {
		body := &def.SetSceneRequest{}
		if err := decode(body); err != nil {
			return nil, err
		}

		if err := body.Validate(); err != nil {
			return nil, err
		}

		return h.SetScene(ctx, body)
	})

}
