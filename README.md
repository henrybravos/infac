# INFAC - Sistema de Facturación Electrónica para Perú

Un servicio API desarrollado en Go para la emisión de comprobantes electrónicos según las normas SUNAT de Perú.

## Características

- ✅ Emisión de Facturas Electrónicas
- ✅ Emisión de Boletas de Venta Electrónicas  
- ✅ Emisión de Notas de Crédito Electrónicas
- ✅ Emisión de Notas de Débito Electrónicas
- 🚧 Anulación de Comprobantes (en desarrollo)
- ✅ Integración con SUNAT via SOAP
- ✅ Soporte para OSE (Operadores de Servicios Electrónicos)
- ✅ Generación de XML en formato UBL 2.1
- ✅ API REST para integración con frontends

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
  password: "MODDATOS"

issuer:
  document_type: "6"
  document_number: "20100070970"
  name: "MI EMPRESA S.A.C."
  # ... más configuraciones
```

4. Ejecutar el servicio:
```bash
go run cmd/api/main.go
```

## Uso de la API

### Crear y enviar una factura

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "type": "01",
    "serie": "F001",
    "number": "00000001",
    "issue_date": "2025-01-15",
    "currency_code": "PEN",
    "customer": {
      "document_type": "6",
      "document_number": "20100070970",
      "name": "CLIENTE EMPRESA S.A.C."
    },
    "lines": [{
      "quantity": 1,
      "unit_code": "NIU",
      "description": "Servicio de consultoría",
      "unit_price": 100.00,
      "taxes": [{
        "type": "IGV",
        "code": "1000",
        "rate": 18.00
      }]
    }]
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
├── internal/
│   ├── config/        # Configuración
│   ├── handlers/      # Controladores HTTP
│   ├── models/        # Modelos de datos
│   └── services/      # Lógica de negocio
└── pkg/
    ├── soap/          # Cliente SOAP para SUNAT
    ├── ubl/           # Generación de XML UBL 2.1
    └── signature/     # Firma digital
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

**Nota**: Este es un proyecto de demostración. Para uso en producción, asegúrate de implementar todas las validaciones de seguridad y cumplir con los requisitos completos de SUNAT.