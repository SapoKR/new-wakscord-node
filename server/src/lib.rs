use serde::{Deserialize};
use tide::Request;

#[derive(Debug, Deserialize)]
pub struct NodeStatus {
    node_id: u16,
    owner: String
}

#[derive(Debug, Deserialize)]
pub struct ProcessStatus {
    total: u64,
    messages: u64,
    tasks: u64
}

#[derive(Debug, Deserialize)]
pub struct StatusJson {
    info: NodeStatus,
    pending: ProcessStatus,
    processed: u32,
    deleted: u32,
    uptime: u32
}

pub fn handleRoot(mut req: Request<()>) -> tide::Result {
    
    Ok()
}