go build -a -ldflags '-extldflags "-static"' -o dragon.exe

rmdir -Recurse .build

mkdir  ./.build/dragon

move dragon.exe .build/dragon/

mkdir .build/dragon/res/dragon
mkdir .build/dragon/data/dragon

copy -Recurse ./control/admin/* .build/dragon/res/dragon/
copy -Recurse ./control/admin/ui/favicon.ico .build/dragon/res/dragon/

cd .build/dragon/

start dragon.exe
