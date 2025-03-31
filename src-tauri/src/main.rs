#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use tauri::{Manager, Window};
use std::sync::Mutex;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

// Define calculation history entry structure
#[derive(Clone, Serialize, Deserialize)]
struct CalculationEntry {
    expression: String,
    result: String,
    timestamp: u64,
}

// State structure to store calculation history
struct CalculatorState {
    history: Mutex<Vec<CalculationEntry>>,
}

#[tauri::command]
fn calculate(expression: String) -> Result<String, String> {
    // This is a stub for the actual calculation
    // In a real implementation, this would parse and compute the expression safely
    // For now, we'll return a placeholder to show the structure
    Ok(format!("42.0"))
}

#[tauri::command]
async fn fetch_currency_rates(base: String) -> Result<HashMap<String, f64>, String> {
    // This would normally make an HTTP request to a currency API
    // For demo purposes, we're returning mock data
    let mut rates = HashMap::new();
    rates.insert("USD".to_string(), 1.0);
    rates.insert("EUR".to_string(), 0.93);
    rates.insert("GBP".to_string(), 0.79);
    rates.insert("JPY".to_string(), 151.15);
    
    Ok(rates)
}

fn main() {
    tauri::Builder::default()
        .manage(CalculatorState {
            history: Mutex::new(Vec::new()),
        })
        .invoke_handler(tauri::generate_handler![
            calculate,
            fetch_currency_rates,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
} 