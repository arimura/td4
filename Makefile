LOGISIM_JAR=logisim-evolution-3.9.0-all.jar
LOGIDIM_URL=https://github.com/logisim-evolution/logisim-evolution/releases/download/v3.9.0/logisim-evolution-3.9.0-all.jar

install: bin/$(LOGISIM_JAR)

bin/$(LOGISIM_JAR):
	mkdir -p bin
	curl -L $(LOGIDIM_URL) -o bin/$(LOGISIM_JAR)

run: bin/$(LOGISIM_JAR)
	java -jar bin/$(LOGISIM_JAR)