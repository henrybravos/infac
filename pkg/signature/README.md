# Certificados Digitales para SUNAT

## Archivos en este directorio

### `certificate_fixed.pfx`
- Certificado digital principal para firma de documentos electrónicos
- Formato: PKCS#12 (PFX)
- Contraseña: `20612790168NEOFORCE`
- Estado: Convertido y optimizado para uso con sunatlib

### `password.txt`
- Contiene credenciales para certificados y SUNAT SOL
- **ADVERTENCIA**: Este archivo contiene información sensible

### `temp/`
- Directorio temporal usado por sunatlib para extraer certificados
- Se crea automáticamente durante el proceso de firma
- Contiene archivos PEM temporales extraídos del PFX

## Firma Digital

La firma digital se maneja completamente a través de `sunatlib`:
- El sistema NO usa implementación propia de firma XML
- Toda la lógica de firma está en la librería externa `github.com/henrybravos/sunatlib`
- Solo se requiere el certificado PFX y la contraseña

## Conversión del Certificado

El certificado original `206127901684LNEOFORCE.pfx` (9.4KB) fue convertido a `certificate_fixed.pfx` (3.8KB) para resolver problemas de compatibilidad con la librería de firma digital.

### Proceso de conversión realizado:
1. Extracción del certificado y clave privada del PFX original
2. Recreación del archivo PFX con formato optimizado
3. Validación de compatibilidad con sunatlib

## Seguridad

- Los certificados PFX están excluidos del control de versiones (`.gitignore`)
- Usar permisos restrictivos: `chmod 600 *.pfx`
- No compartir las contraseñas en texto plano en producción
- El directorio `temp/` se limpia automáticamente

## Uso en el código

```go
err := sunatClient.SetCertificateFromPFX(
    "pkg/signature/certificate_fixed.pfx",
    "20612790168NEOFORCE",
    "pkg/signature/temp",
)
```

## Notas importantes

- Solo usar `certificate_fixed.pfx` - el archivo original ya fue eliminado
- El certificado es específico para el RUC 20612790168
- Requerido para firmar todos los documentos electrónicos enviados a SUNAT