# Tareas Pendientes - INFAC

## 🔄 Estado actual del sistema

El sistema actualmente funciona como un **validador y enviador directo** a SUNAT, pero carece de persistencia real para los documentos. Los borradores se crean en memoria y se pierden inmediatamente.

## 🎯 Mejoras requeridas para producción

### 1. Capa de Persistencia de Documentos
**Estado:** Pendiente  
**Prioridad:** Alta  

**Problema actual:**
- `POST /documents` crea documento en memoria y se pierde
- No hay forma de recuperar borradores
- `POST /send` requiere enviar el JSON completo nuevamente

**Solución propuesta:**
Implementar una capa de persistencia con las siguientes opciones:

#### Opción A: Base de datos SQL
```sql
CREATE TABLE documents (
    id VARCHAR(50) PRIMARY KEY,
    serie VARCHAR(10) NOT NULL,
    number VARCHAR(10) NOT NULL,
    type VARCHAR(2) NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### Opción B: Archivos JSON
```
storage/
├── drafts/
│   ├── F001-00001.json
│   ├── F001-00002.json
│   └── ...
└── sent/
    ├── F001-00001.json
    └── ...
```

#### Opción C: Redis/Cache
```go
redis.Set("draft:F001-00001", documentJSON, 24*time.Hour)
redis.Set("sent:F001-00001", documentJSON, 0) // Sin expiración
```

### 2. Nuevos Endpoints de Gestión
**Estado:** Pendiente  
**Prioridad:** Media

#### 2.1 GET /api/v1/documents/:id
**Propósito:** Recuperar un documento por ID
```bash
GET /api/v1/documents/F001-00001
# Response: documento completo con estado actual
```

#### 2.2 PUT /api/v1/documents/:id
**Propósito:** Actualizar un documento en estado draft
```bash
PUT /api/v1/documents/F001-00001
# Body: campos a actualizar
# Restricción: solo documentos en estado 'draft'
```

#### 2.3 GET /api/v1/documents
**Propósito:** Listar documentos con filtros
```bash
GET /api/v1/documents?status=draft&limit=10
GET /api/v1/documents?serie=F001&from_date=2025-09-01
```

#### 2.4 DELETE /api/v1/documents/:id
**Propósito:** Eliminar borradores no enviados
```bash
DELETE /api/v1/documents/F001-00001
# Restricción: solo documentos en estado 'draft'
```

### 3. Modificación del Endpoint de Envío
**Estado:** Pendiente  
**Prioridad:** Alta

**Cambio propuesto:**
```bash
# Actual (recibe documento completo)
POST /api/v1/documents/send
Body: { documento completo... }

# Propuesto (solo necesita ID)
POST /api/v1/documents/F001-00001/send
Body: {} # vacío o parámetros opcionales
```

**Flujo mejorado:**
1. Sistema busca documento por ID en persistencia
2. Valida que esté en estado `draft`
3. Genera XML, firma y envía a SUNAT
4. Actualiza estado a `sent` o `rejected`
5. Guarda respuesta de SUNAT (CDR)

### 4. Interfaz de Repositorio
**Estado:** Pendiente  
**Prioridad:** Alta

**Implementación sugerida:**
```go
type DocumentRepository interface {
    // CRUD básico
    Save(doc *models.Document) error
    FindByID(id string) (*models.Document, error)
    Update(doc *models.Document) error
    Delete(id string) error
    
    // Consultas especializadas
    FindByStatus(status models.DocumentStatus) ([]*models.Document, error)
    FindBySerie(serie string, limit int) ([]*models.Document, error)
    FindByDateRange(from, to time.Time) ([]*models.Document, error)
    
    // Estadísticas
    CountByStatus() (map[models.DocumentStatus]int, error)
    GetRecentDocuments(limit int) ([]*models.Document, error)
}

// Implementaciones concretas
type JSONFileRepository struct { ... }
type PostgreSQLRepository struct { ... }
type RedisRepository struct { ... }
```

### 5. Validaciones y Estados
**Estado:** Pendiente  
**Prioridad:** Media

**Estados de documento:**
- `draft` - Borrador editable
- `pending` - En proceso de envío
- `sent` - Enviado exitosamente
- `accepted` - Aceptado por SUNAT
- `rejected` - Rechazado por SUNAT
- `cancelled` - Anulado

**Transiciones permitidas:**
```
draft → pending → sent → accepted
draft → pending → sent → rejected
accepted → cancelled (proceso de anulación)
```

### 6. Configuración de Persistencia
**Estado:** Pendiente  
**Prioridad:** Media

**Agregar a config.yaml:**
```yaml
storage:
  type: "json" # json, postgresql, redis
  
  # Para JSON files
  json:
    drafts_path: "storage/drafts"
    sent_path: "storage/sent"
    
  # Para PostgreSQL
  postgresql:
    host: "localhost"
    port: 5432
    database: "infac"
    username: "infac_user"
    password: "secret"
    
  # Para Redis
  redis:
    host: "localhost"
    port: 6379
    database: 0
    password: ""
    draft_ttl: "24h"
```

## 🚀 Beneficios de la implementación

1. **Flujo real de negocio:** Crear → Revisar → Modificar → Enviar
2. **Trazabilidad:** Historial completo de documentos
3. **Recuperación:** No se pierden documentos en caso de fallas
4. **Auditoría:** Registro de cambios y estados
5. **Escalabilidad:** Preparado para múltiples usuarios
6. **API RESTful:** Endpoints estándar para integración

## 📊 Estimación de esfuerzo

| Tarea | Complejidad | Tiempo estimado |
|-------|-------------|-----------------|
| Repositorio JSON | Baja | 2-3 horas |
| Nuevos endpoints | Media | 4-6 horas |
| Modificación /send | Baja | 1-2 horas |
| Validaciones estado | Media | 2-3 horas |
| Configuración | Baja | 1 hora |
| **Total** | | **10-15 horas** |

## 🎯 Priorización sugerida

1. **Fase 1:** Repositorio JSON + modificación /send (funcionalidad básica)
2. **Fase 2:** Nuevos endpoints GET/PUT/DELETE (gestión completa)
3. **Fase 3:** Base de datos SQL (escalabilidad)
4. **Fase 4:** Redis/Cache (performance)

---

**Nota:** Estas mejoras transformarán el sistema de un "demo funcional" a una **aplicación de producción completa** para facturación electrónica.