# API Design Document: Fleet Management System

**Phase:** 1.0  
**Purpose:** This document provides a robust blueprint for the API design of a fleet management system. It includes domain models with updated schemas incorporating cascading privileges, API endpoints, user role hierarchies, entity relationships, and authentication mechanisms. The system manages ports, fleets, assets, and related entities for drone operations, with fine-grained access control.

## User Roles and Hierarchy

### Role Hierarchy Explanation

The system implements cascading privileges through a role hierarchy where higher roles inherit all privileges of lower roles. This allows efficient permission checks, ensuring elevated roles can perform subordinate actions without redundancy. The hierarchy is:

- **Viewer**: Can only view records and their children.
- **Editor**: Inherits all Viewer privileges; can additionally edit records and their children.
- **Owner**: Inherits all Editor (and Viewer) privileges; can additionally create and delete records and their children.

Permissions are enforced in application logic: For any endpoint, verify if the user's effective role level meets or exceeds the minimum required.

Privileges cascade hierarchically to child entities (e.g., Port permissions apply to its Fleets and their Assets). A user's role on a parent entity grants at least Viewer access to children, but the effective child role is the minimum of the inherited parent role and any explicit child assignment.

Access control supports:
- **Global roles**: Default per-user.
- **Per-entity overrides**: Explicit assignments via the Permissions table.
- **Inheritance logic**: If no explicit permission, traverse entity hierarchy (e.g., Asset → Fleet → Port) for the highest inheritable role; fall back to global if none found.

Permission checks:
1. Query explicit Permission for user/entity.
2. If absent, recurse up hierarchy for inheritable role.
3. Fall back to user's global role.
4. Compare role level (Viewer=1, Editor=2, Owner=3) against endpoint minimum.

This enables flexible, scalable access (e.g., Owner on a Port inherits Owner on child Fleets unless overridden).

## Domain Reference

### Entity Relationship Diagram

Below is an ASCII diagram illustrating the relationships between the domain entities:

```
+--------+       +---------+
|  User  |<----->| Session |
+--------+ 1-*   +---------+
   |
   v
+------+
| Port |       
+------+       
   ^           
   | 1-*       
   |           
   v           
+------------+ 
|   Fleet    | 
+------------+ 
   |    |      
   | *-*|      
   v    v      
+--------------+       +-----------+
| Asset-Template |<-----| Component |
+--------------+ 1-*   +-----------+
   ^    ^                      
   |    |                      
   |    |                      
1-*|    |                      
   |    |                      
   v    |                      
+--------------+               
| Manufacturer |               
+--------------+               
   ^                           
   | 1-*                       
   v                           
+-------+                      
| Asset |                      
+-------+                      
   ^                           
   | 1-*                       
   v                           
+----------+                   
| Asset-Part|                  
+----------+                   
```

- **User** 1-* **Session**: A user can have multiple sessions.
- **Port** 1-* **Fleet**: A port can have multiple fleets; a fleet belongs to one port.
- **Fleet** *-* **Asset-Template**: A fleet can use multiple asset templates (many-to-many association).
- **Asset-Template** 1-* **Asset**: An asset template can have multiple asset instances.
- **Asset-Template** 1-* **Component**: An asset template can have multiple components (sub-parts).
- **Manufacturer** 1-* **Asset-Template**: A manufacturer can produce multiple asset templates.
- **Manufacturer** 1-* **Component**: A manufacturer can produce multiple components.
- **Asset** 1-* **Asset-Part**: An asset can have multiple parts (instances of template components).
- **Fleet** 1-* **Asset**: A fleet can have multiple assets.

### Database Schemas

#### Roles
```sql
CREATE TABLE Roles (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,  -- e.g., 'Viewer', 'Editor', 'Owner'
    level INT NOT NULL  -- e.g., 1 for Viewer, 2 for Editor, 3 for Owner
);
```

#### Permissions
```sql
CREATE TABLE Permissions (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    entity_type VARCHAR(255) NOT NULL,  -- e.g., 'Port', 'Fleet', 'Asset'
    entity_id VARCHAR(64) NOT NULL,
    role_id VARCHAR(64) NOT NULL,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (role_id) REFERENCES Roles(id),
    UNIQUE (user_id, entity_type, entity_id)  -- One role per user per entity
);
```

#### Manufacturers
```sql
CREATE TABLE Manufacturers (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255),
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL
);
```

#### Users
```sql
CREATE TABLE Users (
    id VARCHAR(64) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id VARCHAR(64),  -- Global default role
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES Roles(id)
);
```

#### Ports
```sql
CREATE TABLE Ports (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location_lat DOUBLE NOT NULL,
    location_lng DOUBLE NOT NULL,
    location_elevation DOUBLE NOT NULL, -- metres
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL
);
```

#### Fleets
```sql
CREATE TABLE Fleets (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    port_id VARCHAR(64) NOT NULL,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (port_id) REFERENCES Ports(id)
);
```

#### AssetTemplates
```sql
CREATE TABLE AssetTemplates (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    manufacturer_id VARCHAR(64) NOT NULL,
    product_weight BIGINT, -- grams
    product_width INT, -- cm
    product_height INT, -- cm
    product_length INT, -- cm
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (manufacturer_id) REFERENCES Manufacturers(id)
);
```

#### Components
```sql
CREATE TABLE Components (
    id VARCHAR(64) PRIMARY KEY,
    template_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    manufacturer_id VARCHAR(64) NOT NULL,
    product_weight BIGINT, -- grams
    product_width INT, -- cm
    product_height INT, -- cm
    product_length INT, -- cm
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (template_id) REFERENCES AssetTemplates(id),
    FOREIGN KEY (manufacturer_id) REFERENCES Manufacturers(id)
);
```

#### FleetAssetTemplates
```sql
CREATE TABLE FleetAssetTemplates (
    fleet_id VARCHAR(64) NOT NULL,
    template_id VARCHAR(64) NOT NULL,
    PRIMARY KEY (fleet_id, template_id),
    FOREIGN KEY (fleet_id) REFERENCES Fleets(id),
    FOREIGN KEY (template_id) REFERENCES AssetTemplates(id)
);
```

#### Assets
```sql
CREATE TABLE Assets (
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
```

#### AssetParts
```sql
CREATE TABLE AssetParts (
    id VARCHAR(64) PRIMARY KEY,
    asset_id VARCHAR(64) NOT NULL,
    component_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    serial_number VARCHAR(255),
    condition VARCHAR(255),
    notes TEXT,
    FOREIGN KEY (asset_id) REFERENCES Assets(id),
    FOREIGN KEY (component_id) REFERENCES Components(id)
);
```

#### Sessions
```sql
CREATE TABLE Sessions (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    expiry BIGINT NOT NULL,
    time_created BIGINT NOT NULL,
    time_updated BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);
```

## API Reference

### Authentication and Authorization

- All endpoints (except login) require `Authorization: Bearer <session_id>` header.
- Session validation: Check expiry and user existence.
- Role enforcement: Compute effective role per request; return 403 if insufficient.

### User Management

#### /user (POST)
Creates a new user.
```
{
    "username": "newuser",
    "password": "password123",
    "role": "viewer"
}
```
Minimum role: Owner  

#### /user/:user_id (GET)
Returns user data.
```
{
    "id": "980144e66414d1b752ed4e8c3159876ebcc623e611b97d05be8b099518ff08be",
    "type": "user",
    "fields": {
        "username": "newuser",
        "role": "viewer"
    }
}
```
Minimum role: Viewer  

#### /user/:user_id (PATCH)
Updates user fields.
```
{
    "role": "editor"
}
```
Minimum role: Owner  

### Session Management

#### /session (POST)
Creates a session (login).
```
{
    "username": "newuser",
    "password": "password123"
}
```
Minimum role: None  
```
{
    "id": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "user_id": "980144e66414d1b752ed4e8c3159876ebcc623e611b97d05be8b099518ff08be",
    "expiry": 1722276000
}
```

#### /session/:session_id (DELETE)
Deletes the session (logout).  
Minimum role: Viewer  

### Port Management

Ports are physical locations used by one or more fleets.

#### /port (POST)
Creates a new port.
```
{
    "name": "london_biggin_hill",
    "location_lat": 1.23423,
    "location_lng": 50.23234,
    "location_elevation": 123.7
}
```
Minimum role: Owner  

#### /port/:port_id/fleet/:fleet_id (PUT)
Adds fleet to port.  
Minimum role: Editor  

#### /port/:port_id/fleet/:fleet_id (DELETE)
Removes fleet from port.  
Minimum role: Owner  

#### /port/:port_id (GET)
Returns port data.
```
{
    "id": "980144e66414d1b752ed4e8c3159876ebcc623e611b97d05be8b099518ff08be",
    "type": "port",
    "fields": {
        "name": "london_biggin_hill",
        "location_lat": 1.23423,
        "location_lng": 50.23234,
        "location_elevation": 123.7
    },
    "fleet": [
        "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    ]
}
```
Minimum role: Viewer  

### Asset Templates

An asset-template is referenced by an asset.

#### /asset-template (POST)
Creates a new asset template.
```
{
    "name": "hyperdrone_xf-11",
    "manufacturer_id": "some_manufacturer_id",
    "product_width": 80,
    "product_height": 40,
    "product_length": 60
}
```
Minimum role: Owner  

#### /asset-template/:template_id (PATCH)
Updates template fields.
```
{
    "product_weight": 8500
}
```
Minimum role: Editor  

#### /asset-template/:template_id (GET)
Returns template data.
```
{
    "id": "080144e66414d1b752ed4e8c3159876ebcc623e611b97d05be8b099518ff08be",
    "type": "template",
    "fields": {
        "name": "hyperdrone_xf-11",
        "manufacturer_id": "some_manufacturer_id",
        "product_width": 80,
        "product_height": 40,
        "product_length": 60
    },
    "components": [
        "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
        "8e64849654372d661f6dcd75982223e5eea274c42dd53249513a9eb2d2b27980",
        "93aaf4aa03375a0c5b3b3ec3535d82ceb10279e2fb54968e269705cb67724f5a",
        "b6f6116d7812445d7e48730178b3168674a28e0dbf938a4e81825fdd478b5a26"
    ]
}
```
Minimum role: Viewer  

#### /asset-template/:template_id/component (POST)
Creates a component.
```
{
    "name": "Front-left rotor blades",
    "manufacturer_id": "some_manufacturer_id",
    "product_weight": 100,
    "product_width": 80,
    "product_height": 40,
    "product_length": 60
}
```
Minimum role: Owner  

#### /asset-template/:template_id/component/:component_id (PATCH)
Patches component.
```
{
    "product_weight": 85
}
```
Minimum role: Editor  

#### /asset-template/:template_id/component/:component_id (DELETE)
Removes component.  
Minimum role: Owner  

### Fleet Management

#### /fleet (POST)
Creates a new fleet.
```
{
    "name": "Test Fleet A",
    "description": "A fleet for testing the API"
}
```
Minimum role: Owner  

#### /fleet/:fleet_id/template/:template_id (PUT)
Adds template to fleet.  
Minimum role: Editor  

#### /fleet/:fleet_id/template/:template_id (DELETE)
Removes template from fleet.  
Minimum role: Owner  

#### /fleet/:fleet_id/asset/:asset_id (PUT)
Moves asset to fleet.  
Minimum role: Editor  

#### /fleet/:fleet_id/asset/:asset_id (DELETE)
Removes asset from fleet.  
Minimum role: Owner  

#### /fleet/:fleet_id (GET)
Returns fleet data.
```
{
    "id": "980144e66414d1b752ed4e8c3159876ebcc623e611b97d05be8b099518ff08be",
    "type": "fleet",
    "fields": {
        "name": "Test Fleet A",
        "description": "A fleet for testing the API",
        "port": "980144e66414d1b752ed4e8c3159876ebcc623e611b97d05be8b099518ff08be"
    }
}
```
Minimum role: Viewer  

### Asset Instances

An asset instance is an instance of an asset-template that can be managed individually.

#### /asset (POST)
Creates a new asset.
```
{
    "name": "My New Drone A",
    "template": "hyperdrone_xf-11",
    "date_buy": 73287028340,
    "date_install": 697697669,
    "warranty": "1y"
}
```
Minimum role: Owner  

#### /asset/:asset_id (GET)
Returns asset data.
```
{
    "id": "980144e66414d1b752ed4e8c3159876ebcc623e611b97d05be8b099518ff08be",
    "type": "asset",
    "fields": {
        "name": "My New Drone A",
        "template": "hyperdrone_xf-11",
        "date_buy": 73287028340,
        "date_install": 697697669,
        "warranty": "1y"
    },
    "components": [
        "5f873f7e0bd101ba77dbc40a6ea76771cfe5a5e0ef3da03938a282821ef4c0d6",
        "d1a6c49d4c0cb7b729ed20a20f17f4d7c5f3a2d52246d49786019b9fd729e254",
        "e62b95cc85caf807581870b03be48547b5f0edb07160eb10b5d2b320f6a8f49c"
    ]
}
```
Minimum role: Viewer  

#### /asset/:asset_id/inspections/schedule (POST)
Schedules inspection.
```
{
    "timestamp": 7978894354358
}
```
Minimum role: Editor  

#### /asset/:asset_id/inspections/log (POST)
Logs inspection updates.
```
[
    {
        "component_id": "1bfeada3385656172e88398d453d8e22f661f489b92bdde63f076e30ff46099f",
        "serial_number": "656172e88398d4",
        "action": "replaced",
        "condition": "new",
        "notes": "was damaged",
        "time": 78269786307
    }
]
```
Minimum role: Editor