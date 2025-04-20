CREATE TABLE IF NOT EXISTS solar_system_commodity_markets (
    ID uuid,
    Base_Price Double Precision,
    Demand_Quantity Integer,
    Commodity_ID uuid,
    Solar_System_ID uuid,
    Created_At TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    Updated_At TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (ID),
    UNIQUE (Commodity_ID, Solar_System_ID)
);

ALTER TABLE solar_system_commodity_markets ADD CONSTRAINT fk_commodity_id FOREIGN KEY (Commodity_ID) REFERENCES commodities(ID);
ALTER TABLE solar_system_commodity_markets ADD CONSTRAINT fk_solar_system_id FOREIGN KEY (Solar_System_ID) REFERENCES solar_systems(ID);
