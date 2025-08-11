use lambda_runtime::{Error, LambdaEvent};
use serde_json::Value;

pub(crate) async fn function_handler(event: LambdaEvent<Value>) -> Result<String, Error> {
    let payload = event.payload;
    let retval = format!("Payload: {:?}", payload);

    Ok(retval)
}

#[cfg(test)]
mod tests {
    use super::*;
    use lambda_runtime::{Context, LambdaEvent};

    #[tokio::test]
    async fn test_event_handler() {
        let payload = serde_json::from_str(r#"{"foo":"bar"}"#).unwrap();
        let event = LambdaEvent::new(payload, Context::default());
        let response = function_handler(event).await.unwrap();
        assert_eq!(r#"Payload: Object {"foo": String("bar")}"#, response);
    }
}
