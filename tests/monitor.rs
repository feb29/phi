#[macro_use]
extern crate log;
extern crate pretty_env_logger;
extern crate phi;

#[test]
fn monitor() {
    let _ = pretty_env_logger::init();
    let mut mon = phi::Monitor::new(100, 0.8, 0.8);
    mon.push(80.0);
    mon.push(80.0);
    mon.push(70.0);
    mon.push(80.0);
    mon.push(90.0);
    mon.push(200.0);
    mon.push(190.0);
    mon.push(180.0);
    mon.push(280.0);
    mon.push(350.0);

    debug!("{:?}", mon.estimate(500.0));
}
