
pub use config::initialize;
use tide::log;

#[async_std::main]
async fn main() -> tide::Result<()> {
    log::start();

    let config = initialize();
    let value = match config {
        Ok(value) => value,
        Err(err) => {
            panic!("Error: {}", err);
        }
    };

    let mut app = tide::new()
        .at("/").post(handleRoot)
        .listen(
        "127.0.0.1:8080"
        ).await?;

    Ok(())
}