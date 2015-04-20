from distutils.core import setup, Extension

readsymtab = Extension('readsymtab', sources = ['readsymtab.c'])

setup (name = 'readsymtab',
        version = '1.0',
        description = 'Read the symbol table',
        ext_modules = [readsymtab]
)
