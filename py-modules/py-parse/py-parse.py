import ast
import sys


def parse_file(filename):
    with open(filename) as f:
        return ast.parse(f.read(), filename=filename)


def print_body(node):
    try:
        for stmt in node.body:
            print(f'\t {stmt.lineno}', stmt)
            print_body(stmt)
    except AttributeError:
        pass


if __name__ == '__main__':
    tree = parse_file(sys.argv[1])

    for item in ast.walk(tree):
        if isinstance(item, (ast.FunctionDef, ast.AsyncFunctionDef)):
            args = item.args
            args_count = len(args.posonlyargs) + \
                len(args.args) + len(args.kwonlyargs)
            print(item.name, ':', item.lineno, f'({args_count})')
            print_body(item)
