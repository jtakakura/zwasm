// SIMD benchmark: BLAKE2b hashing throughput
// Tests blake2b_simd crate (pure Rust with dynamic SIMD detection)
// Build: cargo build --target wasm32-wasip1 --release

use blake2b_simd::Params;

fn main() {
    let data_size: usize = 1024 * 1024; // 1MB block
    let iterations: usize = 100;

    // Create 1MB of deterministic data
    let mut data = vec![0u8; data_size];
    for (i, byte) in data.iter_mut().enumerate() {
        *byte = (i % 251) as u8;
    }

    let mut final_hash = [0u8; 32];

    for _ in 0..iterations {
        let hash = Params::new()
            .hash_length(32)
            .hash(&data);
        final_hash.copy_from_slice(hash.as_bytes());
        // Feed hash back as prefix to create chain dependency
        data[0] = final_hash[0];
        data[1] = final_hash[1];
    }

    let total_mb = (data_size * iterations) as f64 / (1024.0 * 1024.0);
    println!(
        "blake2b: {:.0} MB hashed, final: {:02x}{:02x}{:02x}{:02x}",
        total_mb, final_hash[0], final_hash[1], final_hash[2], final_hash[3]
    );
}
