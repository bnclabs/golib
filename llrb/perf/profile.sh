count=1000000
go build || exit 1
./perf -count $count -pprof llrb.initial.pprof -mprof llrb.initial.mprof
#./perf -algo mvcc -ops min -count $count -pprof llrb.min.pprof -mprof llrb.min.mprof
#./perf -algo mvcc -ops max -count $count -pprof llrb.max.prof -mprof llrb.max.mprof
#./perf -algo mvcc -ops get -count $count -pprof llrb.get.pprof -mprof llrb.get.mprof
#./perf -algo mvcc -ops upsert -count $count -pprof llrb.upsert.pprof -mprof llrb.upsert.mprof
#./perf -algo mvcc -ops range -count $count -pprof llrb.range.pprof -mprof llrb.upsert.mprof
#./perf -algo mvcc -ops delete -count $count -pprof llrb.del.pprof -mprof llrb.del.mprof

files=initial # min max get upsert range del
for file in $files; do
    file="llrb.${file}.pprof"
    go tool pprof --svg perf $file > $file.svg
done
for file in $files; do
    file="llrb.${file}.mprof"
    go tool pprof --svg --inuse_space perf $file > $file.inuse_space.svg
    go tool pprof --svg --alloc_space perf $file > $file.alloc_space.svg
done
