package x52

import (
	"fmt"
)

// String returns a string representation of the LED
func (led LED) String() string {
	switch led {
	case LedFire:
		return "Fire"

	case LedA:
		return "A"

	case LedB:
		return "B"

	case LedD:
		return "D"

	case LedE:
		return "E"

	case LedT1:
		return "T1"

	case LedT2:
		return "T2"

	case LedT3:
		return "T3"

	case LedPOV:
		return "POV"

	case LedClutch:
		return "Clutch"

	case LedThrottle:
		return "Throttle"
	}

	return fmt.Sprintf("LED(%d)", led)
}

// String returns a string representation of the LED state
func (state LedState) String() string {
	switch state {
	case LedOff:
		return "Off"

	case LedOn:
		return "On"

	case LedRed:
		return "Red"

	case LedAmber:
		return "Amber"

	case LedGreen:
		return "Green"
	}

	return fmt.Sprintf("State(%d)", state)
}

// SetLed sets the state of the given LED. Not all LEDs support all states,
// LedFire and LedThrottle only support LedOn and LedOff states,
// the remaining LEDs support every state except LedOn.
// **Limitation**: This function will not work on a non-pro X52 at this time.
func (ctx *Context) SetLed(led LED, state LedState) error {
	// Make sure that this is a supported device
	// The non-pro X52 doesn't support setting LED states
	if !ctx.HasFeature(FeatureLED) {
		return errNotSupported("setting LED state")
	}

	switch led {
	case LedFire, LedThrottle:
		if state == LedOff {
			bitClear(&ctx.ledMask, uint32(led))
			bitSet(&ctx.updateMask, uint32(led))
		} else if state == LedOn {
			bitSet(&ctx.ledMask, uint32(led))
			bitSet(&ctx.updateMask, uint32(led))
		} else {
			return errNotSupported("invalid state for on/off LED")
		}

	case LedA, LedB, LedD, LedE, LedT1, LedT2, LedT3, LedPOV, LedClutch:
		ledID := uint32(led)
		switch state {
		case LedOff:
			bitClear(&ctx.ledMask, ledID+0)
			bitClear(&ctx.ledMask, ledID+1)

		case LedRed:
			bitSet(&ctx.ledMask, ledID+0)
			bitClear(&ctx.ledMask, ledID+1)

		case LedAmber:
			bitSet(&ctx.ledMask, ledID+0)
			bitSet(&ctx.ledMask, ledID+1)

		case LedGreen:
			bitClear(&ctx.ledMask, ledID+0)
			bitSet(&ctx.ledMask, ledID+1)

		default:
			return errNotSupported("invalid state for color LED")
		}

		bitSet(&ctx.updateMask, ledID+0)
		bitSet(&ctx.updateMask, ledID+1)

	default:
		return errNotSupported("invalid LED identifier")
	}

	return nil
}
