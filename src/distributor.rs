use async_trait::async_trait;

// Contains information we care about in a Query Result
// Not all of the fields in the struct below can be supplied
// by every distributor, as such as wrap with a Option
pub struct ProductInfo {
    pub part_number: String,
    pub manafacturer_part_number: Option<String>,
    pub manafacturer: Option<String>,
    pub description: Option<String>,
    pub unit_price: Option<f64>,
    pub stock: Option<u32>,
    pub product_url: Option<String>,
    pub datasheet_url: Option<String>,
}

// Different Distributors have different ways to initialize
// This handles the configuration for a variety of them
pub enum DistributorConfig {
    ApiKey {
        api_key: String,
    },
    OAuth {
        client_id: String,
        client_secret: String,
    },
}

// A trait that any added Distributor must implement
#[async_trait]
pub trait Distributor {
    // Initialization method for setting up the distributor
    fn initialize(&mut self, config: DistributorConfig);

    // Queries product information by a part number
    async fn query_product_info(&self, part_number: &str) -> Option<ProductInfo>;
}
