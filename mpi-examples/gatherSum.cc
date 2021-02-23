
#include <mpi.h>
#include <iostream>


void sendi(int data, int dest) {
  MPI_Send(&data, 1, MPI_INT, dest, 0, MPI_COMM_WORLD);
}

void recvi(int *data, int srce) {
  MPI_Recv(data, 1, MPI_INT, srce, MPI_ANY_TAG, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
}

void initsr(int* ip, int* np, int* ac, char*** av) {
  MPI_Init(ac, av);
  MPI_Comm_rank(MPI_COMM_WORLD, ip);
  MPI_Comm_size(MPI_COMM_WORLD, np);
}
void donesr(void) {
  MPI_Finalize();
}

int main(int argc, char** argv) {
  
  int proc_id;
  int num_procs;
  initsr(&proc_id, &num_procs, &argc,&argv);

  int N = num_procs;
  if (argc >= 2) {
    N = std::atoi(argv[1]);
  }
  
  int  I = proc_id;
  int  P = num_procs;
  int  n = N/P;

  //
  // Acquire and transform my data.
  int* data = new int[n];
  for (int i=0; i<n; i++) {
    data[i] = (I*n+i);
  }

  //
  // Sum my data.
  int sum = 0;
  for (int i = 0; i < n; i++) {
    sum += data[i];
  }

  //
  // Collectively sum.
  if (I == 0) {
    
    // Gather the partial sums from each of the others.
    for (int from_id=1; from_id < P; from_id++) {
      int sum_i;
      recvi(&sum_i, from_id);
      std::cout << "Received " << sum_i << " from processor " << from_id << "." << std::endl;
      sum += sum_i;
    }
    // Report the sum.
    std::cout << "Ran using " << P << " processors." << std::endl;
    std::cout << "The computed sum is " << sum << "." << std::endl;
    std::cout << "It should be " << ((N >> 1) * (N-1)) << "." << std::endl;
    
  } else {

    // Send to the leader.
    sendi(sum,0);
    
  }

  donesr();
}
