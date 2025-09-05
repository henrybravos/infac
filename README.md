# INFAC - Sistema de Facturación Electrónica para Perú

Un servicio API desarrollado en Go para la emisión de comprobantes electrónicos según las normas SUNAT de Perú, cumpliendo con la Resolución 000193-2020 y el estándar UBL 2.1.

## Características

- ✅ Emisión de Facturas Electrónicas (01)
- ✅ Emisión de Boletas de Venta Electrónicas (03)  
- ✅ Emisión de Notas de Crédito Electrónicas (07)
- ✅ Emisión de Notas de Débito Electrónicas (08)
- ✅ Integración completa con SUNAT via SOAP
- ✅ Firma digital de documentos con certificados PFX
- ✅ Generación de XML en formato UBL 2.1 compliant
- ✅ API REST con validación completa
- ✅ Manejo de términos de pago (Contado/Crédito)
- ✅ Soporte completo para catálogos SUNAT
- 🚧 Anulación de Comprobantes (en desarrollo)

## Requisitos

- Go 1.21+
- Certificado digital válido para firma electrónica (para producción)
- Credenciales SUNAT o OSE

## Instalación

1. Clonar el repositorio:
```bash
git clone <repository-url>
cd infac
```

2. Instalar dependencias:
```bash
go mod download
```

3. Configurar el archivo `config.yaml`:
```yaml
server:
  host: "0.0.0.0"
  port: "8080"

sunat:
  url: "https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService"
  username: "MODDATOS"
  password: "moddatos"

issuer:
  document_type: "6"
  document_number: "20612790168"
  name: "NEOFORCE BUSINESS SOLUTIONS S.A.C."
  trade_name: "NEOFORCE"
  address: "AV. EJEMPLO 123"
  district: "LIMA"
  province: "LIMA"
  department: "LIMA"
  country: "PE"
  email: "contacto@neoforce.pe"
  phone: "+51-1-4251234"
```

4. Configurar certificado digital (ver `pkg/signature/README.md`)

5. Ejecutar el servicio:
```bash
# Desarrollo con hot reload
./scripts/dev.sh

# O directamente
go run cmd/api/main.go
```

## Uso de la API

### Crear una factura al contado

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "serie": "F001",
    "number": "00001",
    "type": "01",
    "issue_date": "2025-09-05",
    "currency_code": "PEN",
    "customer": {
      "document_number": "20123456789",
      "document_type": "6",
      "name": "EMPRESA CLIENTE SAC"
    },
    "payment_terms": {
      "payment_means_code": "Contado"
    },
    "lines": [{
      "quantity": 2.0,
      "unit_code": "NIU",
      "description": "Laptop HP Pavilion",
      "unit_price": 2500.0,
      "taxes": [{
        "type": "IGV",
        "code": "1000",
        "rate": 18.0
      }],
      "product_code": "LAPTOP001"
    }]
  }'
```

### Crear una factura al crédito

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "serie": "F001",
    "number": "00002",
    "type": "01",
    "issue_date": "2025-09-05",
    "due_date": "2025-10-05",
    "currency_code": "PEN",
    "customer": {
      "document_number": "12345678",
      "document_type": "1",
      "name": "JUAN CARLOS PEREZ"
    },
    "payment_terms": {
      "payment_means_code": "Credito",
      "due_date": "2025-10-05T00:00:00Z",
      "amount": 590.0
    },
    "lines": [{
      "quantity": 1.0,
      "unit_code": "NIU",
      "description": "Servicio de consultoria",
      "unit_price": 500.0,
      "taxes": [{
        "type": "IGV",
        "code": "1000",
        "rate": 18.0
      }],
      "product_code": "SERV001"
    }]
  }'
```

### Enviar documento a SUNAT

```bash
curl -X POST http://localhost:8080/api/v1/documents/send \
  -H "Content-Type: application/json" \
  -d '{
    "id": "F001-00001",
    "serie": "F001",
    "number": "00001",
    "type": "01",
    "issue_date": "2025-09-05T00:00:00Z",
    "currency_code": "PEN",
    "issuer": {
      "document_type": "6",
      "document_number": "20612790168",
      "name": "NEOFORCE BUSINESS SOLUTIONS S.A.C.",
      "trade_name": "NEOFORCE"
    },
    "customer": {
      "document_type": "6",
      "document_number": "20123456789",
      "name": "EMPRESA CLIENTE SAC"
    },
    "payment_terms": {
      "payment_means_code": "Contado"
    },
    "lines": [{
      "quantity": 2,
      "unit_code": "NIU",
      "description": "Laptop HP Pavilion",
      "unit_price": 2500,
      "total_price": 5000,
      "taxable_amount": 5000,
      "taxes": [{
        "type": "IGV",
        "code": "1000",
        "rate": 18,
        "amount": 900
      }]
    }],
    "sub_total": 5000,
    "total_taxes": 900,
    "total_amount": 5900
  }'
```

### Consultar estado de resumen diario (para boletas)

```bash
curl http://localhost:8080/api/v1/documents/status/{ticket}
```

### Crear nota de débito

```bash
curl -X POST http://localhost:8883/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "type": "08",
    "serie": "FD01",
    "number": "00000001",
    "issue_date": "2025-01-16",
    "currency_code": "PEN",
    "customer": {
      "document_type": "6",
      "document_number": "20123456789",
      "name": "EMPRESA CLIENTE S.A.C."
    },
    "related_documents": [{
      "document_type": "01",
      "serie": "F001",
      "number": "00000001"
    }],
    "lines": [{
      "quantity": 1,
      "unit_code": "NIU",
      "description": "Intereses por pago tardío",
      "unit_price": 50.00,
      "taxes": [{"type": "IGV", "code": "1000", "rate": 18.00}]
    }]
  }'
```

### Anular un comprobante

```bash
curl -X POST http://localhost:8883/api/v1/documents/void \
  -H "Content-Type: application/json" \
  -d '{
    "document_type": "01",
    "serie": "F001", 
    "number": "00000001",
    "void_date": "2025-01-15",
    "reason": "Error en emision"
  }'
```

## Estructura del Proyecto

```
infac/
├── cmd/api/           # Punto de entrada de la aplicación
├── docs/              # Documentación API con ejemplos JSON
├── internal/
│   ├── config/        # Configuración
│   ├── handlers/      # Controladores HTTP
│   ├── models/        # Modelos de datos y requests
│   └── services/      # Lógica de negocio
├── pkg/
│   ├── ubl/           # Generación de XML UBL 2.1
│   └── signature/     # Firma digital y certificados
└── scripts/           # Scripts de desarrollo (hot reload)
```

## Tipos de Documentos Soportados

| Código | Tipo de Documento | Estado |
|--------|-------------------|---------|
| 01 | Factura | ✅ |
| 03 | Boleta de Venta | ✅ |
| 07 | Nota de Crédito | ✅ |
| 08 | Nota de Débito | ✅ |

## Configuración de Ambientes

### Testing (Beta SUNAT)
```yaml
sunat:
  url: "https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService"
  username: "MODDATOS"
  password: "MODDATOS"
```

### Producción
```yaml
sunat:
  url: "https://e-factura.sunat.gob.pe/ol-ti-itcpfegem/billService"
  username: "TU_RUC + TU_USUARIO"
  password: "TU_CLAVE_SOL"
```

### Con OSE
```yaml
sunat:
  ose:
    enabled: true
    provider: "nubefact"
    url: "https://demo-ose.nubefact.com/ol-ti-itcpfegem/billService"
    username: "tu_usuario_ose"
    password: "tu_clave_ose"
```

## Consideraciones de Seguridad

- Los certificados digitales deben almacenarse de forma segura
- Usar HTTPS en producción
- Las credenciales deben manejarse como variables de entorno
- Implementar rate limiting para la API

## Contribuir

1. Fork el proyecto
2. Crear una rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit los cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## Soporte

Para reportar bugs o solicitar nuevas características, por favor crear un issue en el repositorio.

---

**Nota**: Este es un proyecto de demostración. Para uso en producción, asegúrate de implementar todas las validaciones de seguridad y cumplir con los requisitos completos de SUNAT.# infac
