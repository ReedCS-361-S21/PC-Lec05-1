
#include <mpi.h>
#include <iostream>
#include <iomanip>
#include <cstdlib>

void sendi(int data, int dest) {
  MPI_Send(&data, 1, MPI_INT, dest, 0, MPI_COMM_WORLD);
}

void recvi(int *data, int srce) {
  MPI_Recv(data, 1, MPI_INT, srce, MPI_ANY_TAG, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
}


void sendiv(int* data, int n, int dest) {
  MPI_Send(data, n, MPI_INT, dest, 0, MPI_COMM_WORLD);
}

void recviv(int *data, int n, int srce) {
  MPI_Recv(data, n, MPI_INT, srce, MPI_ANY_TAG, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
}

void initsr(int* ip, int* np, int* ac, char*** av) {
  MPI_Init(ac, av);
  MPI_Comm_rank(MPI_COMM_WORLD, ip);
  MPI_Comm_size(MPI_COMM_WORLD, np);
}

void donesr(void) {
  MPI_Finalize();
}

int randint(int max) {
  return rand() % max;
}

void swap(int* x, int* y) {
  int tmp = *x;
  *x = *y;
  *y = tmp;
}

void minleft(int* data, int n) {
  for (int i=1; i<n; i++) {
    if (data[i] < data[0]) {
      swap(data,data+i);
    }
  }
}

void maxright(int* data, int n) {
  for (int i=0; i<n-1; i++) {
    if (data[i] > data[n-1]) {
      swap(data+n-1,data+i);
    }
  }
}

void report(/* int r, */ int* data, int n, int I, int P) {
  //
  // Report your values.
  for (int i=0; i<n; i++) {
    // std::cout << std::setfill('0') << std::setw(4) << r << ". ";
    std::cout << std::setfill('0') << std::setw(4) << I << ": x[";
    std::cout << std::setfill('0') << std::setw(4) << n*I+i << "] = ";
    std::cout << std::setfill('0') << std::setw(4) << data[i] << std::endl;
  }
}

void oesort(int* data, int n, int I, int P) {
  for (int round = 0; round < n*P; round++) {
    if (((round + I) % 2) == 0) {
      // Look right. Send first.
      if (I < P-1) {
	int value;
	maxright(data,n);
	sendi(data[n-1],I+1);
	recvi(&value,I+1);
	if (value < data[n-1]) {
	  data[n-1] = value;
	}
      }
    } else {
      // Look left. Recv first.
      if (I > 0) {
	int value;
	minleft(data,n);
	recvi(&value,I-1);
	sendi(data[0],I-1);
	if (value > data[0]) {
	  data[0] = value;
	}
      }
    }
  }
} 

void sort(int* data, int n) {
  for (int i=1; i<n; i++) {
    int x = data[i];
    int j = i;
    while ((j > 0) && (data[j-1] > x)) {
      data[j] = data[j-1];
      j--;
    }
    data[j] =  x;
  }
}

int main(int argc, char** argv) {

     
  int proc_id;
  int num_procs;
  initsr(&proc_id, &num_procs, &argc,&argv);

  int  N = num_procs;
  
  if (argc >= 2) {
    N = std::atoi(argv[1]);
  }
  

  int  I = proc_id;
  int  P = num_procs;
  int  n = N/P;

  srand(time(NULL) + I);
  //
  // Acquire and transform my data.
  int* data = new int[n];
  for (int i=0; i<n; i++) {
    data[i] = randint(N);
  }

  oesort(data,n,I,P);
  sort(data,n);

  if (I == 0) {
    report(data,n,0,P);
    for (int id = 1; id < P; id++) {
      recviv(data,n,id);
      report(data,n,id,P);
    }
  } else {
    sendiv(data,n,0);
  }
  
  donesr();
}
