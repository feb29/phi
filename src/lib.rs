// Package phi implements `The Phi Accrual Failure Detector` defined in:
//    http://ddg.jaist.ac.jp/pub/HDY+04.pdf

mod monitor;
pub use monitor::Monitor;

const E3: f64 = 1e3;
const E6: f64 = 1e6;

const DEFAULT_THRESHOLD: f64 = 16.0;
