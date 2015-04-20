#include <Python.h>
#include "Python-ast.h"
#include "code.h"
#include "symtable.h"
#include "structmember.h"

static PyObject*
readsymtab(PyObject* self, PyObject* args)
{
    PyObject* obj;
 
    if (!PyArg_ParseTuple(args, "O", &obj))
        return NULL;

    if (Py_TYPE(obj) != &PySTEntry_Type) {
        fprintf(stderr, "Ooops wrong type\n");
    }

    PySTEntryObject *st = (PySTEntryObject *)obj;
 
    return Py_BuildValue("(iiiiiii)",
                         st->ste_free,
                         st->ste_child_free,
                         st->ste_generator,
                         st->ste_varargs,
                         st->ste_varkeywords,
                         st->ste_returns_value,
                         st->ste_needs_class_closure);
}
 
static PyMethodDef readsymtabmethods[] =
{
     {"readsymtab", readsymtab, METH_VARARGS, "Read the symbol table."},
     {NULL, NULL, 0, NULL}
};

static struct PyModuleDef readsymtabmodule = {
    PyModuleDef_HEAD_INIT,
    "readsymtab",
    "Read the symbol table.",
    -1,
    readsymtabmethods,
    NULL,
    NULL,
    NULL,
    NULL
};

 
PyMODINIT_FUNC
PyInit_readsymtab(void)
{
    return PyModule_Create(&readsymtabmodule);
}
