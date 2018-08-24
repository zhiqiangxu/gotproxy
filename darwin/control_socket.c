#include "control_socket.h"

int Connect() {
    struct ctl_info ctl_info;
    struct sockaddr_ctl sc;

    int sock = socket(PF_SYSTEM, SOCK_DGRAM, SYSPROTO_CONTROL);
    if (sock < 0) {
        perror("new socket");
        return -1;
    }

    bzero(&ctl_info, sizeof(ctl_info));
    strcpy(ctl_info.ctl_name, MYBUNDLEID);

    if (ioctl(sock, CTLIOCGINFO, &ctl_info) == -1) {
        perror("ioctl");
        return -1;
    }

    bzero(&sc, sizeof(sc));
    sc.sc_len = sizeof(sc);
    sc.sc_family = AF_SYSTEM;
    sc.ss_sysaddr = SYSPROTO_CONTROL;
    sc.sc_id = ctl_info.ctl_id;
    sc.sc_unit = 0;

    if (connect(sock, (struct sockaddr*)&sc, sizeof(sc))) {
        perror("connect");
        return -1;
    }

    return sock;
}

bool StartRedirect(int sock, uint16_t port) {
    struct tproxy_user_param user_param = {port};
    if (setsockopt(sock, SYSPROTO_CONTROL, TPROXY_ON, &user_param, sizeof(user_param)) == -1) return false;

    return true;

}

bool StopClose(int sock) {

    if (setsockopt(sock, SYSPROTO_CONTROL, TPROXY_OFF, NULL, 0) == -1) return false;

    return close(sock) == 0;
}