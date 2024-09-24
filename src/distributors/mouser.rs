use crate::distributor::{Distributor, DistributorConfig, ProductInfo};
use async_trait::async_trait;
use reqwest::Client;
use serde::{Deserialize, Serialize};

// ========== Mouser API Response ==========
#[derive(Debug, Deserialize)]
pub struct MouserApiResponse {
    #[serde(rename = "Errors")]
    pub errors: Option<Vec<Error>>,

    #[serde(rename = "SearchResults")]
    pub search_results: Option<SearchResults>,
}

#[derive(Debug, Deserialize)]
pub struct Error {
    #[serde(rename = "Id")]
    pub id: i32,

    #[serde(rename = "Code")]
    pub code: String,

    #[serde(rename = "Message")]
    pub message: String,

    #[serde(rename = "ResourceKey")]
    pub resource_key: String,

    #[serde(rename = "ResourceFormatString")]
    pub resource_format_string: String,

    #[serde(rename = "ResourceFormatString2")]
    pub resource_format_string2: String,

    #[serde(rename = "PropertyName")]
    pub property_name: String,
}

#[derive(Debug, Deserialize)]
pub struct SearchResults {
    #[serde(rename = "NumberOfResult")]
    pub number_of_result: i32,

    #[serde(rename = "Parts")]
    pub parts: Vec<Part>,
}

#[derive(Debug, Deserialize)]
pub struct Part {
    #[serde(rename = "MouserPartNumber")]
    pub part_number: String,

    #[serde(rename = "ManufacturerPartNumber")]
    pub manufacturer_part_number: String,

    #[serde(rename = "Manufacturer")]
    pub manufacturer: String,

    #[serde(rename = "Description")]
    pub description: String,

    #[serde(rename = "PriceBreaks")]
    pub price_breaks: Vec<PriceBreak>,

    #[serde(rename = "AvailabilityInStock")]
    pub availability_in_stock: String,

    #[serde(rename = "unitprice")]
    pub unit_price: f64,

    #[serde(rename = "ProductDetailUrl")]
    pub product_url: String,

    #[serde(rename = "DataSheetUrl")]
    pub data_sheet_url: String,
}

#[derive(Debug, Deserialize)]
pub struct PriceBreak {
    #[serde(rename = "Quantity")]
    pub quantity: i32,

    #[serde(rename = "Price")]
    pub price: String,

    #[serde(rename = "Currency")]
    pub currency: String,
}

// Define the request body structure for the POST request.
#[derive(Serialize, Debug)]
struct MouserRequestBody {
    #[serde(rename = "SearchByPartRequest")]
    search_by_part_request: SearchByPartRequest,
}

#[derive(Serialize, Debug)]
struct SearchByPartRequest {
    #[serde(rename = "mouserPartNumbers")]
    part_numbers: Vec<String>,

    #[serde(rename = "apiKey")]
    api_key: String,
}

// ====================

pub struct MouserDistributor {
    api_key: String,
    client: Client,
    base_url: String,
}

impl MouserDistributor {
    pub fn new() -> Self {
        Self {
            api_key: String::new(),
            client: Client::new(),
            base_url: String::new(),
        }
    }
    async fn perform_query_mouser(
        &self,
        part_number: &str,
    ) -> Result<MouserApiResponse, reqwest::Error> {
        let url = format!(
            "{}/search/partnumber?apiKey={}",
            self.base_url, self.api_key
        );
        // Create the request body using the part number and API key
        let request_body = MouserRequestBody {
            search_by_part_request: SearchByPartRequest {
                part_numbers: vec![part_number.to_string()],
                api_key: self.api_key.clone(),
            },
        };
        let response = self.client.post(url).json(&request_body).send().await?;
        response.json::<MouserApiResponse>().await
    }
}

#[async_trait]
impl Distributor for MouserDistributor {
    fn initialize(&mut self, config: DistributorConfig) {
        if let DistributorConfig::ApiKey { api_key } = config {
            self.api_key = api_key;
            self.base_url = "https://api.mouser.com/api/v1/".to_string();
        }
    }

    async fn query_product_info(&self, part_number: &str) -> Option<ProductInfo> {
        match self.perform_query_mouser(part_number).await {
            Ok(response) => {
                if let Some(search_results) = response.search_results {
                    if let Some(part) = search_results.parts.get(0) {
                        return Some(ProductInfo {
                            part_number: part.part_number.clone(),
                            manafacturer_part_number: Some(part.manufacturer_part_number.clone()),
                            description: Some(part.description.clone()),
                            manafacturer: Some(part.manufacturer.clone()),
                            unit_price: Some(part.unit_price),
                            stock: Some(part.availability_in_stock.parse().unwrap_or(0)),
                            datasheet_url: Some(part.data_sheet_url.clone()),
                            product_url: Some(part.product_url.clone()),
                        });
                    }
                }
                None
            }
            Err(_) => None,
        }
    }
}
