use miller_rabin;
use std::io;
use std::io::Read;
fn main() {
    let mut out = String::new();
    let mut buf = String::new();
    let _ = io::stdin().read_to_string(&mut buf);
    for st in buf.split_whitespace() {
        let num: u128 = st.parse().unwrap();
        if miller_rabin::is_prime(&num, 200) {
            out += &format!("{} is prime.\n", &num);
        } else {
            out += &format!("{} is not prime.\n", &num);
        }
    }
    println!("{}", out);
}
