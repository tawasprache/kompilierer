# tawa

tawa ist eine deutschsprachige Programmieresprache.

warum? das grund ist mir egal.

## ebnf-grammatik

```
Datei = "paket" <string> (<string> "ist" "importiert")* Typdeklaration* Funktion* .
Typdeklaration = "typ" <ident> "ist" Art .
Art = Struktur | Zeiger | <ident> .
Struktur = "struktur" "(" Strukturfield* ")" .
Strukturfield = <ident> Art .
Zeiger = "zeiger" "auf" Art .
Funktion = "funk" <ident> "(" (Funktionsargument ("," Funktionsargument)*)? ")" (":" Art)? Expression .
Funktionsargument = <ident> ":" Art .
Expression = (Bedingung | Definierung | Zuweisung | Funktionsaufruf | Logik | Cast | Integer | Löschen | Neu | Stack | Dereferenzierung | <ident> | Block) Postfix*? .
Bedingung = "wenn" Expression Expression ("sonst" Expression)? .
Definierung = <ident> ":" Art? "=" Expression .
Zuweisung = <ident> "=" Expression .
Funktionsaufruf = <ident> "(" (Expression ("," Expression)*)? ")" .
Logik = ("Wahr" | "Falsch") .
Cast = "cast" Expression "nach" Art .
Integer = <int> .
Löschen = "lösche" Expression .
Neu = "neu" Expression .
Stack = "stack" Strukturinitialisierung .
Strukturinitialisierung = <ident> "(" Strukturinitialisierungsfield* ")" .
Strukturinitialisierungsfield = <ident> "ist" Expression .
Dereferenzierung = "deref" Expression .
Block = ("(" Expression* ")") .
Postfix = Fieldoperator | Zuweisungsoperator .
Fieldoperator = "." <ident> .
Zuweisungsoperator = "=" Expression .
```

## expressionen

### bedingung

```
wenn <expr> <expr> (sonst <expr>)
```

Dieses sind Konditionelle Expressionen. Sie erlaubt man, um Entscheidungen zu machen.
Die erste Expression müsste ein Logikwert sein.

### definierung

```
<name>: (art) = <expr>
```

Diese erlauben man, um Variablen zu machen.

### zuweisung

```
<name> = <expr>
```

### funktionsaufruf

```
<name>(<expr>, <expr>, ...)
```

### logik

```
Falsch
```

oder

```
Wahr
```

Du weißt diese bereits.

### cast

```
cast <expr> nach <art>
```

### integer

```
[0-9]*
```

Du weißt diese bereits.

### löschen

```
lösche <expr>
```

Löscht die allokiertes speicher.

### neu

```
neu <expr>
```

Kopiert ein Wert von Stack-Speicher nach Heap-Speicher.

Müsst löscht sein.

### dereferenzierung

```
deref <expr>
```

Kopiert ein Wert von Heap-Speicher nach Stack-Speicher.

### Variable

```
<name>
```

Der Wert der benannten Variable.

### Block

```
(
	<expr>
	<expr>
	...
)
```

Viele expressionen. Der Wert ist gleich auf der Wert letztes Expression.

### Fieldexpressionen

```
<expr>.<name>
```

Zugriff auf Fields des eine Wert.

## Arten

### Primitiven

- `nichts`: nichts
- `ganz`: eine ganzes zahl
- `ganz8`
- `ganz16`
- `ganz32`
- `ganz64`
- `vzlosganz`: eine vorzeichenlose ganzes zahl
- `vzlosg8`
- `vzlosg16`
- `vzlosg32`
- `vzlosg64`

### Zeigern

```
zeiger auf <art>
```

### Strukturen

```
struktur <name> (
	<name> <art>
	(<name> <art>)
	(<name> <art>)
	(...)
)
```
