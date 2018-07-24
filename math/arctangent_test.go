package math

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/MeetFrankie/decimal"
)

func TestAtan(t *testing.T) {
	const N = 100
	diff := decimal.WithPrecision(N)
	eps := decimal.New(1, N)
	for i, tt := range [...]struct {
		x, r string
	}{
		0: {"0", "0"},
		1: {".500", "0.4636476090008061162142562314612144020285370542861202638109330887201978641657417053006002839848878926"},
		2: {"-.500", "-0.4636476090008061162142562314612144020285370542861202638109330887201978641657417053006002839848878926"},
		3: {".9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", "0.7853981633974483096156608458198757210492923498437764552437361480769541015715522496570087063355292669"},
		4: {"-.9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", "-0.7853981633974483096156608458198757210492923498437764552437361480769541015715522496570087063355292669"},
		5: {"1.00", "0.7853981633974483096156608458198757210492923498437764552437361480769541015715522496570087063355292670"},
		6: {"-1.00", "-0.7853981633974483096156608458198757210492923498437764552437361480769541015715522496570087063355292670"},
		7: {".999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", "0.785398163397448309615660845819875721049292349843776455243736148076954101571552249657008706335529266995537021628320576661773461152387645557931339852032120279362571025675484630276389"},
		8: {"100.0", "1.560796660108231381024981575430471893537215347143176270859532877957451649939045719334570767484384444"},
		9: {"0." + strings.Repeat("9", 1000), "0.7853981633974483096156608458198757210492923498437764552437361480769541015715522496570087063355292669955370216283205766617734611523876455579313398520321202793625710256754846302763899111557372387325954911072027439164833615321189120584466957913178004772864121417308650871526135816620533484018150622853184311467516515788970437203802302407073135229288410919731475900028326326372051166303460367379853779023582643175914398979882730465293454831529482762796370186155949906873918379714381812228069845457529872824584183406101641607715053487365988061842976755449652359256926348042940732941880961687046169173512830001420317863158902069464428356894474022934092946803671102253062383575366373963427626980699223147308855049890280322554902160086045399534074436928274901296768028374999995932445124877649329332040240796487561148638367270756606305770633361712588154827970427525007844596882216468833020953551542944172868258995633726071888671827898907159705884468984379894454644451330428067016532504819691527989773041050497"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			z := decimal.WithPrecision(N)
			x, _ := new(decimal.Big).SetString(tt.x)
			r, _ := new(decimal.Big).SetString(tt.r)

			Atan(z, x)
			if z.Cmp(r) != 0 && diff.Sub(r, z).CmpAbs(eps) > 0 {
				t.Errorf(`#%d: Atan(%s)
wanted: %s
got   : %s
diff  : %s
`, i, x, r, z, diff)
			}
		})
	}
}

var atan_X, _ = new(decimal.Big).SetString("0.7853981633974483096156608458198757210492923498437764552437361480769541015715522496570087063355292670")

func BenchmarkAtan(b *testing.B) {
	for _, prec := range benchPrecs {
		b.Run(fmt.Sprintf("%d", prec), func(b *testing.B) {
			z := decimal.WithPrecision(prec)
			for j := 0; j < b.N; j++ {
				Atan(z, atan_X)
			}
			gB = z
		})
	}
}

func TestAtan2(t *testing.T) {
	const N = 1000
	for i, tt := range [...]struct {
		y, x, r string
	}{
		// Atan2(NaN, NaN) = NaN
		0: {"NaN", "NaN", "NaN"},

		// Atan2(y, NaN) = NaN
		1: {"NaN", "NaN", "NaN"},
		2: {"NaN", "NaN", "NaN"},

		// Atan2(NaN, x) = NaN
		3: {"NaN", "0", "NaN"},
		4: {"NaN", "1", "NaN"},

		// Atan2(+/-0, x>=0) = +/-0
		5: {"0", "0", "0"},
		6: {"-0", "0", "-0"},
		7: {"0", "1", "0"},
		8: {"-0", "1", "-0"},

		// Atan2(+/-0, x<=-0) = +/-pi
		9:  {"0", "-0", pos(_pi, N)},
		10: {"-0", "-0", neg(_pi, N)},
		11: {"0", "-1", pos(_pi, N)},
		12: {"-0", "-1", neg(_pi, N)},

		// Atan2(y>0, 0) = +pi/2
		13: {"1", "0", pos(_pi_2, N)},

		// Atan2(y<0, 0) = -pi/2
		14: {"-1", "0", neg(_pi_2, N)},

		// Atan2(+/-Inf, +Inf) = +/-pi/4
		15: {"+Inf", "+Inf", pos(_pi_4, N)},
		16: {"-Inf", "+Inf", neg(_pi_4, N)},

		// Atan2(+/-Inf, -Inf) = +/-3pi/4
		17: {"+Inf", "-Inf", pos(_3pi_4, N)},
		18: {"-Inf", "-Inf", neg(_3pi_4, N)},

		// Atan2(y, +Inf) = 0
		19: {"-1", "+Inf", "0"},
		20: {"0", "+Inf", "0"},
		21: {"1", "+Inf", "0"},

		// Atan2(y>0, -Inf) = +pi
		22: {"1", "-Inf", pos(_pi, N)},

		// Atan2(y<0, -Inf) = -pi
		23: {"-1", "-Inf", neg(_pi, N)},

		// Atan2(+/-Inf, x) = +/-pi/2
		24: {"+Inf", "-1", pos(_pi_2, N)},
		25: {"-Inf", "-1", neg(_pi_2, N)},
		26: {"+Inf", "0", pos(_pi_2, N)},
		27: {"-Inf", "0", neg(_pi_2, N)},
		28: {"+Inf", "1", pos(_pi_2, N)},
		29: {"-Inf", "1", neg(_pi_2, N)},

		// Atan2(y,x>0) = Atan(y/x)
		30: {"-1", "1", neg(_pi_4, N)},
		31: {"1", "1", pos(_pi_4, N)},

		// Atan2(y>=0, x<0) = Atan(y/x) + pi
		32: {"0", "-1", pos(_pi, N)},
		33: {"1", "-1", pos(_3pi_4, N)},

		// Atan2(y<0, x<0) = Atan(y/x) - pi
		34: {"-1", "-1", neg(_3pi_4, N)},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			z := decimal.WithPrecision(N)
			y, _ := new(decimal.Big).SetString(tt.y)
			x, _ := new(decimal.Big).SetString(tt.x)
			r, _ := new(decimal.Big).SetString(tt.r)
			Atan2(z, y, x)
			if z.Cmp(r) != 0 {
				t.Fatalf(`#%d: Atan2(%s, %s)
wanted: %s
got   : %s
`, i, y, x, r, z)
			}
		})
	}
}
