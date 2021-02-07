<?php

class cls {
	public $prop1;

	public function __construct() {
		$this->prop1 = "ppp";
	}

	public function funccccc() {
		echo 'test';
	}
}

$b = true;
$b1 = false;

$i = 123;

$f = 3.1415;

$s = 'a ba a ba, "(){}?:"><?!@#$%^&*-=_+|\'\\中国';

$n = null;

$o = new cls();

$a = [
	$b,
	3.14 => $b1,
	$i,
	$f,
	'str' => $s,
	'null' => $n,
    'obj' => $o,
];

$filter = [
	$b,
	$b1,
    $i,
    $f,
	$s,
	$n,
	$o,
	$a,
];

if (@$_GET['q']!='json') {
	foreach ($filter as $key => $value) {
		echo serialize($value);
		if ($key != count($filter)-1) {
			echo "\n";
		}
	}
}
else {
	foreach ($filter as $key => $value) {
		echo json_encode($value, JSON_UNESCAPED_UNICODE|JSON_HEX_TAG|JSON_HEX_AMP|JSON_FORCE_OBJECT);
		if ($key != count($filter)-1) {
			echo "\n";
		}
	}
}