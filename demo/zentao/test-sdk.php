#!/usr/bin/env php
<?php
/**
[case]
title=Use Zentao SDK to do interface testing
cid=0
pid=0

[group]
  1. check first product name >> `.+`

[esac]
*/
include_once 'vendor/sdk.php';

$zentao      = new \zentao();

$extraFields = array('title', 'products', 'productStats');    // 自定义返回字段
$result      = $zentao->getProductList(array(), $extraFields);
$result      = json_decode($result);
$name = $result->result->productStats[0]->name;

print(">> $name\n");
