int exported_global;
double exported_global_2;
extern int imported_global;

void imported(void);

void exported(void) {
	exported_global = 5;
	exported_global_2 = 3.14;
	imported_global = 3;
	imported();
}

double exported_2(void) {
	return exported_global_2;
}
