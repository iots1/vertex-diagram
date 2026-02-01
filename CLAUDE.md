# Vertex-Diagram Backend - Project Overview & Integration Guide

> **Purpose**: MongoDB-backed REST API to replace IndexedDB local storage for ChartDB frontend. All users view shared diagram data without local persistence.

---

## 📋 Project Summary

**Type**: Go REST API Backend
**Framework**: Fiber v2 (Express-like)
**Database**: MongoDB (Document-based)
**Architecture**: Clean Architecture with Dependency Injection
**Current Status**: Functional but lacks security features

---

## 🏗️ Architecture Overview

```
Domain Layer (Business Rules)
    ↓ (implements)
Usecase Layer (Business Logic)
    ↓ (uses)
Repository Layer (MongoDB CRUD)
    ↓ (reads/writes)
Delivery Layer (HTTP Handlers)
    ↓ (serves)
Frontend (ChartDB)
```

---

## 🔌 API Endpoints

### Base URL
```
http://localhost:8080/api
```

### Diagram Endpoints
| Method | Endpoint | Purpose |
|--------|----------|---------|
| `GET` | `/diagrams` | List all diagrams (metadata only, no content) |
| `GET` | `/diagrams/:id` | Get single diagram with full content + tables + relationships |
| `POST` | `/diagrams` | Create or update diagram |
| `DELETE` | `/diagrams/:id` | Delete diagram and cascade-delete related data |

### Config Endpoints
| Method | Endpoint | Purpose |
|--------|----------|---------|
| `GET` | `/config` | Get global configuration |
| `POST` | `/config` | Update global configuration |

---

## 💾 Database Schema

### Collections

#### `diagrams`
```json
{
  "_id": "string/ObjectID",
  "name": "string",
  "content": {
    // Large JSON blob containing full diagram structure
    "tables": [...],
    "relationships": [...],
    // ... other diagram data
  },
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```
**Index**: `_id` (primary)
**Note**: Content excluded from list queries for performance

#### `tables`
```json
{
  "_id": "ObjectID",
  "diagram_id": "string/ObjectID",
  "table_id": "string",
  "name": "string",
  "schema": "string",
  "fields": [{...}],
  "indexes": [{...}],
  "color": "string",
  "x": int,
  "y": int,
  "is_view": bool,
  "order": int,
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```
**Indexes**: `diagram_id` (for fast filtering)

#### `relationships`
```json
{
  "_id": "ObjectID",
  "diagram_id": "string/ObjectID",
  "relationship_id": "string",
  "name": "string",
  "source_table_id": "string",
  "target_table_id": "string",
  "source_field_id": "string",
  "target_field_id": "string",
  "type": "string",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```
**Indexes**: `diagram_id` (for fast filtering)

#### `config`
```json
{
  "_id": "global",
  "default_diagram_id": "string/ObjectID"
}
```
**Special**: Only one document (singleton pattern)

---

## 🔄 Data Flow

### Saving a Diagram
```
POST /api/diagrams
  ↓
Extract from diagram.content:
  - tables array → Insert/update in tables collection
  - relationships array → Insert/update in relationships collection
  ↓
Store diagram document with updated timestamp
  ↓
Response: Updated diagram object
```

### Loading a Diagram
```
GET /api/diagrams/:id
  ↓
Fetch diagram from diagrams collection
Fetch all tables WHERE diagram_id = :id
Fetch all relationships WHERE diagram_id = :id
  ↓
Merge tables + relationships back into diagram.content
  ↓
Response: Complete diagram with nested data
```

### Deleting a Diagram
```
DELETE /api/diagrams/:id
  ↓
Cascade delete:
  1. Delete all tables WHERE diagram_id = :id
  2. Delete all relationships WHERE diagram_id = :id
  3. Delete diagram document
  ↓
Response: 204 No Content
```

---

## ⚙️ Configuration

### Environment Variables
```env
# MongoDB connection
MONGO_URI=mongodb://root:u2035@49.0.124.7:27037

# Server
PORT=8080

# Database
DB_NAME=vertex_db
```

### Server Settings
- **Body Limit**: 50MB (supports large diagram payloads)
- **Request Timeout**: 5 seconds
- **CORS**: Allows all origins (`*`)
- **Hot Reload**: Enabled via Air tool

---

## 🚨 Known Issues & Security Concerns

### Critical Issues
1. **❌ No Authentication** - API is completely open, no user verification
2. **❌ No Authorization** - No access control, all users see all diagrams
3. **❌ No Input Validation** - Request payloads not validated before storage
4. **❌ No Rate Limiting** - API vulnerable to abuse/DoS attacks
5. **❌ CORS Too Permissive** - Allows requests from any origin

### Code Quality Issues
1. **🐛 Bug in config_usecase.go:31** - Incorrect parameter order in Upsert call
2. **⚠️ No Logging** - Minimal error logging for debugging
3. **⚠️ No API Documentation** - No Swagger/OpenAPI specs
4. **⚠️ Type Conversion** - JSON unmarshaling converts numbers to float64, requires conversion

### Data Isolation Issues
1. **🔓 Single-Tenant Design** - No multi-tenancy support
2. **🔓 Shared Data Store** - All diagrams visible to all users
3. **🔓 No Row-Level Security** - No permission checks per diagram

---

## 🔐 Recommended Security Enhancements

### Phase 1: Basic Authentication (MVP)
```go
// Add JWT middleware
// Add user claims to context
// Validate token on protected endpoints
```

### Phase 2: Authorization
```go
// Add user ID to diagrams collection
// Filter diagrams by user_id
// Implement role-based access control
```

### Phase 3: Hardening
```go
// Add input validation (use pkg/validator)
// Implement rate limiting (use middleware)
// Add comprehensive logging (use zap/logrus)
// Implement soft deletes for audit trail
```

---

## 🔗 Integration with ChartDB Frontend

### Frontend → Backend Flow

1. **Load Dashboard** (list all diagrams)
   ```typescript
   // chartdb frontend
   const response = await fetch('http://localhost:8080/api/diagrams');
   const diagrams = await response.json();
   // Display: [{id, name, updatedAt, createdAt}, ...]
   ```

2. **Open Diagram for Editing**
   ```typescript
   const response = await fetch(`http://localhost:8080/api/diagrams/${diagramId}`);
   const diagram = await response.json();
   // Render: {id, name, content: {tables, relationships}, ...}
   // NO local storage - ALL data from server
   ```

3. **Save Diagram Changes**
   ```typescript
   await fetch(`http://localhost:8080/api/diagrams`, {
     method: 'POST',
     body: JSON.stringify({
       id: diagramId,
       name: diagramName,
       content: {
         tables: [...],
         relationships: [...],
         // ... all diagram data
       }
     })
   });
   // Response: Updated diagram object
   ```

4. **Delete Diagram**
   ```typescript
   await fetch(`http://localhost:8080/api/diagrams/${diagramId}`, {
     method: 'DELETE'
   });
   // Response: 204 No Content
   ```

5. **Get/Set Default Diagram**
   ```typescript
   // On app load
   const config = await fetch('http://localhost:8080/api/config');
   const defaultId = (await config.json()).default_diagram_id;

   // On user action
   await fetch('http://localhost:8080/api/config', {
     method: 'POST',
     body: JSON.stringify({ default_diagram_id: newDiagramId })
   });
   ```

### Key Requirements Met ✓
- ✅ **No Local Storage**: All data fetched from server
- ✅ **Shared Data**: All users view same diagrams
- ✅ **Server-Side Persistence**: MongoDB stores all diagrams
- ✅ **Real-Time Sync**: Users see latest version when diagram is reloaded

---

## 📊 Current Data Sharing Model

**All Users → Single Database → Shared Diagrams**

```
User A ──┐
User B ──┤
User C ──┼──→ MongoDB (vertex_db) ──→ All users see same data
User D ──┤
...      ──┘
```

**Behavior**:
- No user isolation
- No permission checks
- All diagrams visible to everyone
- Changes are immediately visible to all users (after reload)

---

## 🛠️ Development & Deployment

### Local Development
```bash
# Start server with hot reload
air

# Or manual build
go build -o vertex-diagram .
./vertex-diagram
```

### Dependencies
```
- Fiber v2.52.11    (Web framework)
- MongoDB Driver    (Database driver)
- godotenv          (.env file loader)
- UUID              (ID generation)
```

### Build & Deploy
```bash
# Build binary
go build -o vertex-diagram .

# Run
PORT=8080 DB_NAME=vertex_db ./vertex-diagram
```

---

## 📝 File Structure

```
vertex-diagram/
├── main.go                          # Application entry point
├── go.mod / go.sum                  # Dependencies
├── .env                             # Configuration
├── .air.toml                        # Hot reload config
│
├── domain/                          # Business domain
│   ├── diagram.go
│   ├── table.go
│   ├── relationship.go
│   └── config.go
│
├── repository/                      # Data access (MongoDB)
│   ├── mongo_repository.go
│   ├── mongo_table_repository.go
│   ├── mongo_relationship_repository.go
│   └── mongo_config_repository.go
│
├── usecase/                         # Business logic
│   ├── diagram_usecase.go
│   └── config_usecase.go
│
├── delivery/http/                   # HTTP handlers
│   ├── diagram_handler.go
│   └── config_handler.go
│
└── infrastructure/                  # Infrastructure setup
    ├── config/
    │   └── config.go
    └── database/
        └── database.go
```

---

## ✅ Integration Checklist

- [ ] Start vertex-diagram backend server
- [ ] Verify MongoDB connection and indexes created
- [ ] Test all API endpoints with Postman/curl
- [ ] Update ChartDB frontend API base URL to backend
- [ ] Replace IndexedDB queries with HTTP fetch calls
- [ ] Remove localStorage.* calls from ChartDB
- [ ] Implement error handling for network failures
- [ ] Test real-time sync across multiple browser tabs
- [ ] Implement authentication (Phase 2)
- [ ] Add input validation (Phase 2)
- [ ] Deploy to production environment

---

## 🎯 Next Steps

### Immediate (MVP)
1. ✅ Backend API functional
2. ✅ Database schema defined
3. Update frontend to use HTTP API instead of IndexedDB
4. Test end-to-end diagram save/load flow

### Short-term (Security)
1. Implement JWT authentication
2. Add user_id to diagrams collection
3. Filter diagrams by user
4. Add request validation

### Medium-term (Production)
1. Implement rate limiting
2. Add comprehensive logging
3. Set up monitoring/alerts
4. Performance optimization
5. Backup/recovery procedures

---

## 📞 Endpoints Summary

| Feature | Endpoint | Method |
|---------|----------|--------|
| List diagrams | `/api/diagrams` | GET |
| Get diagram | `/api/diagrams/:id` | GET |
| Save diagram | `/api/diagrams` | POST |
| Delete diagram | `/api/diagrams/:id` | DELETE |
| Get config | `/api/config` | GET |
| Save config | `/api/config` | POST |

---

## 🔍 Monitoring

To verify the backend is working:

```bash
# Check if server is running
curl http://localhost:8080/api/diagrams

# Should return: [] (empty array if no diagrams)
```

---

**Last Updated**: 2026-02-01
**Generated by**: Claude Code (Multi-Agent Analysis)
