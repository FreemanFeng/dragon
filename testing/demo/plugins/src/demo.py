import sys
from pydemo import main

if __name__ == '__main__':
    if len(sys.argv) > 1:
        port = sys.argv[1]
        main.main(port)
    else:
        main.main()