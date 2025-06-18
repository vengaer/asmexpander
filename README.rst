AsmExpander
===========

Expand GNU AS-style macros for easier stepping with gdb.

Nested macros **are** supported.

Passing the following snippet

.. code-block::

    .macro isolssbl r0, r1
        leal -1(\r0), \r1
        notl \r1
        andl \r1, \r0
    .endm

    foo:
        movl $0x21, %eax
        isolssbl %eax, %r8d

to stdin yields

.. code-block::

    .macro isolssbl r0, r1
        leal -1(\r0), \r1
        notl \r1
        andl \r1, \r0
    .endm

    foo:
        movl $0x21, %eax
        leal -1(%eax), %r8d
        notl %r8d
        andl %r8d, %eax

