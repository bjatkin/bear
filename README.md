# Bear
Bear is an error package for making errors in go awesome.
It's designed to be extreamly flexible, letting you use errors in the way that best fits your code base.
Bear focuses on making errors easy to interpret especially at scale.
This means Bear focuses on giving you good ways go gropu errors together and generate meaningful reports about the errors your code is generating.

# Roadmap
* update the fmt options so the apply to all errors in the stack.
    Options like FmtNoParent, and FmtNoStack should apply to all errors in the stack.
    In this sense these options govern no only the error itself but also all it's parents.
    This introduces some questiosn such as, should these be different than just ErrOptions?
    Should this option be respected no matter where it occures in the stack, or should the root child alone be the one to dictate?
    How can these options be exposed such that Error() can be easily controlled?