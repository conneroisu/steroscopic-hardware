EXEC_NAME = range_code_demo
SOURCE_FILE_NAMES = range_code demo
OBJS = $(foreach src,$(SOURCE_FILE_NAMES),$(src).o)
CFLAGS += -g

# .o rules are implicit
$(EXEC_NAME) : $(OBJS)
	cc -g -o $@ $^

.PHONY : clean
clean :
	@for i in $(OBJS); do \
		if [ -f $$i ]; then \
			rm $$i; \
		fi \
	done

	@if [ -f $(EXEC_NAME) ]; then \
		rm $(EXEC_NAME); \
	fi