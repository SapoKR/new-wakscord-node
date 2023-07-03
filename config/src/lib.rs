use std::env;

pub struct Config {
    Host: String,
    Port: u16,
    MaxConcurrent: u32,
    WaitConcurrent: u32,
    MessageQueueSize: u32,
    Key: String,
    ID: u8,
    Owner: String,
}

pub fn initialize() -> Result<Config, Box<dyn std::error::Error>> {
    let cfg: Result<Config, Box<dyn std::error::Error>> = Ok(Config{
        Host: env::var("HOST").unwrap_or("0.0.0.0".to_owned()),
        Port: env::var("PORT").unwrap_or("3000".to_owned()).parse::<u16>()?,
        MaxConcurrent: env::var("MAX_CONCURRENT").unwrap_or("500".to_owned()).parse::<u32>()?,
        WaitConcurrent: env::var("WAIT_CONCURRENT").unwrap_or("1".to_owned()).parse::<u32>()?,
        MessageQueueSize: env::var("MESSAGE_QUEUE_SIZE").unwrap_or("100".to_owned()).parse::<u32>()?,
        Key: env::var("KEY").unwrap_or("wakscord".to_owned()),
        ID: env::var("ID").unwrap_or("0".to_owned()).parse::<u8>()?,
        Owner: env::var("OWNER").unwrap_or("Unknown".to_owned()),
    });

    cfg
}