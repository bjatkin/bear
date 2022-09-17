# Bear
Bear is an error package for making errors in go awesome.
It's designed to be extreamly flexible, letting you use errors in the way that best fits your code base.
Bear focuses on making errors easy to interpret especially at scale.
This means Bear focuses on giving you good ways go gropu errors together and generate meaningful reports about the errors your code is generating.

# Roadmap
There are future features that are planned for this library.
They are ordered roughly in the order of priority.
This list is likely to change a lot and none of these features should be considerd gaurenteed.

* Add a Wrap method to the Template type

* Add a one stack option to print only the longest stack (e.g. the stack of the most senior error)

* Improve the error stack. Right now it's including files like proc.go and asm_amd64.s

* Add more options to transform JSON (filters, ordering, extra fields, etc.)

* Options to transform labels (filters, combinations, extra lables, etc.)

* Options to transform tags (filters, combinations, extra tags, ext.)

* Options to transform metrics (filters, combinations, extra metrics, etx.)

* Combining Metric and FMetric into a single type, this will allow them to show up in the same place in the resulting json and prevent any differences showing up between the two types

