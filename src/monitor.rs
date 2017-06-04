// Monitor implements `Holt's Linear Method` (DoubleExponentialSmoothing).
// `Holt's Linear Method` is good for non-seasonal data with a trend.
//
//    Level:    L[i]   = a*X[i] + (1−a)*(L[i-1] + T[i-1])
//    Trend:    T[i]   = b*(L[i] − L[i−1]) + (1−b)*T[i−1]
//    Forecast: F[i+1] = L[i] + T[i]
//
// a and b are factor for smoothing observed signals.
pub struct Monitor {
    cap: usize,
    sum: f64,
    sos: f64, // sum of square

    signals: Vec<f64>,
    lf: f64,
    tf: f64,
    levels: Vec<f64>,
    trends: Vec<f64>,
}

impl Monitor {
    pub fn new(cap: usize, lf: f64, tf: f64) -> Self {
        let signals = Vec::with_capacity(cap);
        let sum = 0.0;
        let sos = 0.0;
        let lf = bounded(lf);
        let tf = bounded(tf);
        let levels = Vec::new();
        let trends = Vec::new();
        Monitor {
            cap,
            sum,
            sos,
            signals,
            lf,
            tf,
            levels,
            trends,
        }
    }

    fn level(&self, i: usize, x: f64) -> f64 {
        if i == 0 {
            x
        } else {
            self.lf * x + (1. - self.lf) * (self.levels[i - 1] + self.trends[i - 1])
        }
    }

    fn trend(&self, i: usize, x: f64) -> f64 {
        if i == 0 {
            0.
        } else {
            self.tf * (x - self.levels[i - 1]) + (1. - self.tf) * self.trends[i - 1]
        }
    }

    pub fn push(&mut self, x: f64) {
        self.signals.push(x);
        let i = self.signals.len() - 1;
        let x = self.pushn(i, x);
        self.sum += x;
        self.sos += x * x;
    }

    fn pushn(&mut self, i: usize, x: f64) -> f64 {
        let l1 = self.level(i, x);
        self.levels.push(l1);
        let l2 = self.levels[i];
        let t = self.trend(i, l2);
        self.trends.push(t);
        l2
    }

    pub fn estimate(&self, d: f64) -> f64 {
        let len = self.levels.len() as f64;
        let mean = self.sum / len;
        let var = self.sos / len - mean * mean;
        let stddev = var.sqrt();
        phi_value(mean, stddev, d)
    }

    fn failure(&self) -> f64 {
        let i = self.signals.len();
        let x = self.levels[i] + self.trends[i];
        unimplemented!()
    }
}

fn phi_value(mean: f64, stddev: f64, d: f64) -> f64 {
    let y = (d - mean) / stddev;
    let p = (-y * (1.5976 + 0.070566 * y * y)).exp();
    if d > mean {
        -(p / (1.0 + p)).log10()
    } else {
        -(1.0 - 1.0 / (1.0 + p)).log10()
    }
}

fn bounded(f: f64) -> f64 {
    if f >= 1.0 {
        0.999
    } else if f <= 0. {
        0.100
    } else {
        f
    }
}
