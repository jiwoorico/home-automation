// Code generated by devicegen. DO NOT EDIT.

package domain

import (
	device "github.com/jakewright/home-automation/libraries/go/device"
	def "github.com/jakewright/home-automation/libraries/go/device/def"
	oops "github.com/jakewright/home-automation/libraries/go/oops"
	ptr "github.com/jakewright/home-automation/libraries/go/ptr"
)

type MegaParProfileProperties struct {
	Brightness *int64
	Power      *bool
	Rgb        *device.RGB
	Strobe     *int64
}

func (p *MegaParProfileProperties) unmarshal(m map[string]interface{}) error {
	for property, value := range m {
		switch property {
		case "brightness":
			f, ok := value.(float64)
			if !ok {
				return oops.BadRequest("property 'brightness' was unexpected type")
			}

			i := int64(f)
			if f != float64(i) {
				return oops.BadRequest("property 'brightness' should be an integer")
			}

			if i < 0 {
				return oops.BadRequest("property 'brightness' should be ≥ 0")
			}

			if i > 255 {
				return oops.BadRequest("property 'brightness' should be ≤ 255")
			}

			p.Brightness = &i

		case "power":
			b, ok := value.(bool)
			if !ok {
				return oops.BadRequest("property 'power' was unexpected type")
			}

			p.Power = &b

		case "rgb":
			s, ok := value.(string)
			if !ok {
				return oops.BadRequest("property 'rgb' was unexpected type")
			}

			rgb := &device.RGB{}
			if err := rgb.UnmarshalText([]byte(s)); err != nil {
				return oops.WithMessage(err, "property 'rgb': failed to parse %q as RGB value", s)
			}

			p.Rgb = rgb

		case "strobe":
			f, ok := value.(float64)
			if !ok {
				return oops.BadRequest("property 'strobe' was unexpected type")
			}

			i := int64(f)
			if f != float64(i) {
				return oops.BadRequest("property 'strobe' should be an integer")
			}

			if i < 0 {
				return oops.BadRequest("property 'strobe' should be ≥ 0")
			}

			if i > 255 {
				return oops.BadRequest("property 'strobe' should be ≤ 255")
			}

			p.Strobe = &i

		default:
			return oops.BadRequest("property %q not known", property)
		}
	}

	return nil
}

func (p *MegaParProfileProperties) describe() map[string]*def.Property {
	// Dereference all of the pointers. This makes it easier
	// to assert the values are as expected in tests.
	var brightness int64
	if p.Brightness != nil {
		brightness = *p.Brightness
	}

	var power bool
	if p.Power != nil {
		power = *p.Power
	}

	var rgb device.RGB
	if p.Rgb != nil {
		rgb = *p.Rgb
	}

	var strobe int64
	if p.Strobe != nil {
		strobe = *p.Strobe
	}

	return map[string]*def.Property{
		"brightness": {
			Value: brightness,
			Type:  "int",
			Min:   ptr.Float64(0),
			Max:   ptr.Float64(255),
		},
		"power": {
			Value: power,
			Type:  "bool",
		},
		"rgb": {
			Value: rgb,
			Type:  "rgb",
		},
		"strobe": {
			Value: strobe,
			Type:  "int",
			Min:   ptr.Float64(0),
			Max:   ptr.Float64(255),
		},
	}
}
