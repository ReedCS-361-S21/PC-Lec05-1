#include <mpi.h>
#include <iostream>
#include <iomanip>
#include <cstring>
#include <chrono>
#include <cstdlib>

#define BIG_ENOUGH_STRING 100

void sends(char* data, int dest) {
  MPI_Send(data, std::strlen(data)+1, MPI_CHAR, dest, 0, MPI_COMM_WORLD);
}

void recvs(char* data, int srce) {
  MPI_Recv(data, BIG_ENOUGH_STRING+1, MPI_CHAR, srce, MPI_ANY_TAG, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
  data[BIG_ENOUGH_STRING] = 0;
}

void initsr(int* ip, int* np, int* ac, char*** av) {
  MPI_Init(ac, av);
  MPI_Comm_rank(MPI_COMM_WORLD, ip);
  MPI_Comm_size(MPI_COMM_WORLD, np);
}
void donesr(void) {
  MPI_Finalize();
}

void barrier(void) {
  MPI_Barrier (MPI_COMM_WORLD);
}

void flipcap(char* x, int p) {
  if (p < std::strlen(x)) {
    if ('a' <= x[p] && x[p] <= 'z') {
      x[p] -= 32;
    } else if ('A' <= x[p] && x[p] <= 'Z') {
      x[p] += 32;
    }
  }
}

int randint(int max) {
  return rand() % max;
}

int main(int argc, char** argv) {

  int proc_id;
  int num_procs;
  initsr(&proc_id, &num_procs, &argc,&argv);

  int T = 10;
  if (argc >= 2) {
    T = std::atoi(argv[1]);
  }
  int len = 5;
  if (argc >= 3) {
    len = std::strlen(argv[2]);
  }

  int  I = proc_id;
  int  P = num_procs;

  srand(time(NULL) + I);
  char* message = new char[len+1];
  if (I == 0) {
    if (I == 0 && argc >= 3) {
      std::strncpy(message,argv[2],len+1);
    } else {
      std::strncpy(message,"hello",len+1);
    }
  }

  int Ipred = (I + P - 1) % P;
  int Isucc = (I + 1) % P;
  
  if (I == 0) {
    sends(message,Isucc);
  }
  for (int t=1; t<=T; t++) {
    recvs(message,Ipred);

    std::cout << std::setw(3) << std::setfill('0') << t << ". ";
    std::cout << std::setw(2) << std::setfill('0') << I << ": ";
    std::cout << "Received '" << message << "'." << std::endl;
    if (t < T || I != 0) {
      flipcap(message,randint(len));
      std::cout << std::setw(3) << std::setfill('0') << t << ". ";
      std::cout << std::setw(2) << std::setfill('0') << I << ": ";
      std::cout << "Sent out '" << message << "'." << std::endl;
      sends(message,Isucc);
    }
  }

  barrier();
  donesr();
}
