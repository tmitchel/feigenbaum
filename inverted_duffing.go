package main 

////////////////////////////////////////////////////
// Purpose: To solve the Inverted Duffing         //
// Oscillator given set of initial conditions and //
// parameter values                               //
// Return: A single pdf containing a 2D plot of x //
// vs y (Vx)                                      //
////////////////////////////////////////////////////

import (
    // "fmt"
    "flag"
    "math"
    "strconv"
    "image/color"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    )

////////////////////////////////////////////////////
// Purpose: Simple struct to hold a point for     //
// plotting                                       //
// Variables: X, Y coordinates and an int to      //
// signal when a result blows up and can't be     //
// plotted                                        //
////////////////////////////////////////////////////
type point struct {
    x, y float64
    breaker int
}

////////////////////////////////////////////////////
// Purpose: Do the iterative calculations         //
// Returns: A buffered channel packed with point  //
// structs                                        //
////////////////////////////////////////////////////
func iter(x0, y0, F, dt float64) chan point {
    channel := make(chan point, 400)

    t0 := 0.
    // the iteration is done here
    go func () {
        for x, y, t := x0, y0, t0; ; x, y, t = dt*y + x, dt*(F*math.Cos(t)-0.5*y+x-math.Pow(x, 3)) + y, t + dt {
            channel <- point{x: x, y: y}
        }
    } ()
    return channel
}

func main() {

    // Command-line options
    F := flag.Float64("F", 0.24, "Constant F")
    x0 := flag.Float64("x0", 0, "Initial value for x")
    y0 := flag.Float64("y0", 0, "Initial value for y (dx/dt)")
    t := flag.Int("t", 100, "Number of second")
    dt := flag.Int("dt", 1000, "Step Resolution (-dt=10 gives 10 steps per second)")
    max_min_comp := flag.Bool("comp", false, "Compare F=0.24 and F=0.35")
    flag.Parse()

    // If true, plot highest F value vs lowest
    if *max_min_comp {
        // channels holding high/low results
        low := iter(*x0, *y0, 0.24, 1./float64(*dt))
        high := iter(*x0, *y0, 0.35, 1./float64(*dt))

        p, err := plot.New()
        if err != nil {
            panic(err)
        }

        nsteps := (*t) * (*dt)
        points_low := make(plotter.XYs, nsteps)
        points_high := make(plotter.XYs, nsteps)

        // read from channels to fill plots
        for i := 0; i < nsteps; i++ {
            temp_low := <-low
            temp_high := <-high

            points_low[i].X = temp_low.x
            points_low[i].Y = temp_low.y

            points_high[i].X = temp_high.x
            points_high[i].Y = temp_high.y

        }

        // make the plots prettier
        plot_low, _ := plotter.NewLine(points_low)
        plot_low.Color = color.RGBA{R:255}
        plot_high, _ := plotter.NewLine(points_high)
        plot_high.Color = color.RGBA{B:255}
        p.Add(plot_low, plot_high)

        p.X.Label.Text = "X"
        p.Y.Label.Text = "Y=dx/dt"
        p.Title.Text = "Poincare Section F=0.24 vs F=0.35"

        p.Legend.Add("F=0.24", plot_low)
        p.Legend.Add("F=0.35", plot_high)

        // save the pdf
        p.Save(600, 400, "iduff_comp.pdf")

    } else {
        // channel to hold results for user chosen F value
        results := iter(*x0, *y0, *F, 1./float64(*dt))

        p, err := plot.New()
        if err != nil {
            panic(err)
        }

        nsteps := (*t) * (*dt)
        points := make(plotter.XYs, nsteps)

        // read from channel to fill plot
        for i := 0; i < nsteps; i++ {
            temp := <-results

            points[i].X = temp.x
            points[i].Y = temp.y
        }

        plot, _ := plotter.NewLine(points)
        p.Add(plot)

        F_val := strconv.FormatFloat(*F, 'f', -1, 64)

        p.X.Label.Text = "X"
        p.Y.Label.Text = "Y=dx/dt"
        p.Title.Text = "Poincare Section F="+F_val
        p.Save(600, 400, "iduff_F"+F_val+".pdf")

    }
    
}