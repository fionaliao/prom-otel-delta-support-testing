In prometheus repo:

```
git checkout fionaliao/basic-delta-support

make build

 ./prometheus --web.enable-otlp-receiver --enable-feature=otlp-native-delta-ingestion --enable-feature=native-histograms --config.file=config/testdata/otlp_allow_utf8.good.yml
```

In this repo:
```
go run cmd/otel-counter-writer/main.go
```

Demo queries:

http://localhost:9090/query?g0.expr=sum_over_time%28test_counter_total%7Botel_temporality%3D%22DeltaTemporality%22%7D%5B1m%5D%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0&g1.expr=increase%28test_counter_total%7Botel_temporality%3D%22CumulativeTemporality%22%7D%5B1m%5D%29&g1.show_tree=0&g1.tab=graph&g1.range_input=1h&g1.res_type=auto&g1.res_density=medium&g1.display_mode=lines&g1.show_exemplars=0&g2.expr=test_counter_total%7Botel_temporality%3D%22DeltaTemporality%22%7D%5B1m%5D&g2.show_tree=0&g2.tab=table&g2.range_input=1h&g2.res_type=auto&g2.res_density=medium&g2.display_mode=lines&g2.show_exemplars=0&g3.expr=test_counter_total%7Botel_temporality%3D%22CumulativeTemporality%22%7D%5B1m%5D&g3.show_tree=0&g3.tab=table&g3.range_input=1h&g3.res_type=auto&g3.res_density=medium&g3.display_mode=lines&g3.show_exemplars=0&g4.expr=histogram_quantile%280.5%2C+sum%28sum_over_time%28test_histogram%7Botel_temporality%3D%22DeltaTemporality%22%7D%5B1m%5D%29%29%29&g4.show_tree=0&g4.tab=graph&g4.range_input=5m&g4.res_type=auto&g4.res_density=medium&g4.display_mode=lines&g4.show_exemplars=0&g5.expr=histogram_quantile%280.5%2C+sum%28rate%28test_histogram%7Botel_temporality%3D%22CumulativeTemporality%22%7D%5B1m%5D%29%29%29&g5.show_tree=0&g5.tab=graph&g5.range_input=5m&g5.res_type=auto&g5.res_density=medium&g5.display_mode=lines&g5.show_exemplars=0