BIN = test_basic test_batch_write ingest_data
CXXFLAGS = -I ../rocksdb/include/ -std=gnu++17 -L ../rocksdb -lrocksdb -lpthread -lz -lsnappy -lboost_program_options -ldl

ifeq ($(MAKECMDGOALS), debug)
    CXXFLAGS += -g -ggdb
else
    CXXFLAGS += -O3
endif

all: release

debug: $(BIN)

release: $(BIN)

test_basic: test_basic.cpp common.h
	g++ -std=c++17 -o $@ $< $(CXXFLAGS)

test_batch_write: test_batch_write.cpp common.h
	g++ -std=c++17 -o $@ $< $(CXXFLAGS)

ingest_data: ingest_data.cpp common.h
	g++ -std=c++17 -o $@ $< $(CXXFLAGS)

clean:
	rm -f $(BIN)
