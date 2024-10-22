pub trait PrintFullError {
    fn full(&self) -> String;
}

impl PrintFullError for anyhow::Error {
    fn full(&self) -> String {
        format!("{:#}", self)
    }
}
