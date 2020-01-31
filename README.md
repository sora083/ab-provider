# ab-provider
airticket information provider via ABROAD


### how to start to use go module
```
go mod init
go build
```
* [Go Modules](https://qiita.com/propella/items/e49bccc88f3cc2407745)

https://github.com/isucon/isucon8-qualify/blob/master/webapp/go/src/torb/app.go


```
  private String getCity(List<City> cities) {
    return cities.stream()
        .map(city -> String.format("%s (%s)", city.getName(), city.getCode()))
        .collect(Collectors.toList())
        .toString();
  }

  private String getPriceDescription(Price price) {
    return String.format(
        "%s〜%s(commission: %s)",
        df.format(price.getMin()), df.format(price.getMax()), df.format(price.getCommission()));
  }

  private String convertQrImage(String url) {
    Map<EncodeHintType, ErrorCorrectionLevel> hints = new HashMap<>();
    hints.put(EncodeHintType.ERROR_CORRECTION, ErrorCorrectionLevel.M);

    try (ByteArrayOutputStream output = new ByteArrayOutputStream()) {
      QRCodeWriter writer = new QRCodeWriter();
      BitMatrix bitMatrix = writer.encode(url, BarcodeFormat.QR_CODE, QR_SIZE, QR_SIZE, hints);
      MatrixToImageWriter.writeToStream(bitMatrix, "png", output);
      return Base64.encodeBase64String(output.toByteArray());
    } catch (Exception e) {
      log.error("error: {}", e);
    }
    return null;
  }
  ```