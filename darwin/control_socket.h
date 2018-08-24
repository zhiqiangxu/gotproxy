#include <stdio.h>
#include <stdlib.h>
#include <strings.h>
#include <unistd.h>
#include <sys/socket.h>
#include <sys/ioctl.h>
#include <sys/sys_domain.h>
#include <sys/kern_control.h>
#include <stdbool.h>

#define MYBUNDLEID "com.qtt.xuzhiqiang.gotproxy"

#define TPROXY_ON 1
#define TPROXY_OFF 2

struct tproxy_user_param {
    uint16_t port;
};

extern int Connect();

extern bool StartRedirect(int sock, uint16_t port);

extern bool StopClose(int sock);