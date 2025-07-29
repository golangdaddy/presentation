-- Initialize the database with all required tables

-- Create Roles table
CREATE TABLE IF NOT EXISTS Roles (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL
);

-- Insert default roles
INSERT INTO Roles (id, name, level) VALUES 
    ('viewer', 'Viewer', 1),
    ('reporter', 'Reporter', 2),
    ('editor', 'Editor', 3),
    ('owner', 'Owner', 4)
ON CONFLICT (id) DO NOTHING;

-- Create Users table
CREATE TABLE IF NOT EXISTS Users (
    id VARCHAR(64) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    role_id VARCHAR(64) DEFAULT 'viewer',
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES Roles(id)
);

-- Create Sessions table
CREATE TABLE IF NOT EXISTS Sessions (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    expiry BIGINT NOT NULL,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Create Permissions table
CREATE TABLE IF NOT EXISTS Permissions (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    entity_id VARCHAR(64) NOT NULL,
    role_id VARCHAR(64) NOT NULL,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (role_id) REFERENCES Roles(id),
    UNIQUE (user_id, entity_type, entity_id)
);

-- Create Ports table
CREATE TABLE IF NOT EXISTS Ports (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location_lat DOUBLE PRECISION NOT NULL,
    location_lng DOUBLE PRECISION NOT NULL,
    location_elevation DOUBLE PRECISION NOT NULL,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL
);

-- Create Manufacturers table
CREATE TABLE IF NOT EXISTS Manufacturers (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255),
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL
);

-- Create AssetTemplates table
CREATE TABLE IF NOT EXISTS AssetTemplates (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    manufacturer_id VARCHAR(64) NOT NULL,
    product_weight BIGINT,
    product_width INT,
    product_height INT,
    product_length INT,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (manufacturer_id) REFERENCES Manufacturers(id)
);

-- Create Components table
CREATE TABLE IF NOT EXISTS Components (
    id VARCHAR(64) PRIMARY KEY,
    template_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    manufacturer_id VARCHAR(64) NOT NULL,
    product_weight BIGINT,
    product_width INT,
    product_height INT,
    product_length INT,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (template_id) REFERENCES AssetTemplates(id),
    FOREIGN KEY (manufacturer_id) REFERENCES Manufacturers(id)
);

-- Create Fleets table
CREATE TABLE IF NOT EXISTS Fleets (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    port_id VARCHAR(64) NOT NULL,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (port_id) REFERENCES Ports(id)
);

-- Create FleetAssetTemplates table
CREATE TABLE IF NOT EXISTS FleetAssetTemplates (
    fleet_id VARCHAR(64) NOT NULL,
    template_id VARCHAR(64) NOT NULL,
    PRIMARY KEY (fleet_id, template_id),
    FOREIGN KEY (fleet_id) REFERENCES Fleets(id),
    FOREIGN KEY (template_id) REFERENCES AssetTemplates(id)
);

-- Create Assets table
CREATE TABLE IF NOT EXISTS Assets (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    template_id VARCHAR(64) NOT NULL,
    fleet_id VARCHAR(64),
    date_buy BIGINT,
    date_install BIGINT,
    warranty VARCHAR(255),
    FOREIGN KEY (template_id) REFERENCES AssetTemplates(id),
    FOREIGN KEY (fleet_id) REFERENCES Fleets(id)
);

-- Create AssetParts table
CREATE TABLE IF NOT EXISTS AssetParts (
    id VARCHAR(64) PRIMARY KEY,
    asset_id VARCHAR(64) NOT NULL,
    component_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    serial_number VARCHAR(255),
    condition VARCHAR(255),
    notes TEXT,
    inspection_frequency BIGINT,
    FOREIGN KEY (asset_id) REFERENCES Assets(id),
    FOREIGN KEY (component_id) REFERENCES Components(id)
);

-- Create Inspections table
CREATE TABLE IF NOT EXISTS Inspections (
    id VARCHAR(64) PRIMARY KEY,
    asset_id VARCHAR(64) NOT NULL,
    component_id VARCHAR(64),
    serial_number VARCHAR(255),
    action VARCHAR(255),
    condition VARCHAR(255),
    notes TEXT,
    time BIGINT NOT NULL,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (asset_id) REFERENCES Assets(id),
    FOREIGN KEY (component_id) REFERENCES Components(id)
);

-- Create Attachments table
CREATE TABLE IF NOT EXISTS Attachments (
    id VARCHAR(64) PRIMARY KEY,
    entity_type VARCHAR(255) NOT NULL,
    entity_id VARCHAR(64) NOT NULL,
    uri VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    mime_type VARCHAR(255),
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON Sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON Sessions(expiry);
CREATE INDEX IF NOT EXISTS idx_permissions_user_entity ON Permissions(user_id, entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_fleets_port_id ON Fleets(port_id);
CREATE INDEX IF NOT EXISTS idx_assets_fleet_id ON Assets(fleet_id);
CREATE INDEX IF NOT EXISTS idx_assets_template_id ON Assets(template_id);
CREATE INDEX IF NOT EXISTS idx_asset_parts_asset_id ON AssetParts(asset_id);
CREATE INDEX IF NOT EXISTS idx_inspections_asset_id ON Inspections(asset_id);
CREATE INDEX IF NOT EXISTS idx_attachments_entity ON Attachments(entity_type, entity_id); 