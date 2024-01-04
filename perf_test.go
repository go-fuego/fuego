package fuego

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTiming_String(t *testing.T) {
	t.Run("no desc", func(t *testing.T) {
		timing := Timing{
			Name: "test",
			Dur:  time.Duration(100) * time.Millisecond,
		}

		require.Equal(t, "test;dur=100", timing.String())
	})

	t.Run("with desc", func(t *testing.T) {
		timing := Timing{
			Name: "test",
			Dur:  time.Duration(300) * time.Millisecond,
			Desc: "test desc",
		}

		require.Equal(t, "test;dur=300;desc=\"test desc\"", timing.String())
	})
}
