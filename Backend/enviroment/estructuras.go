package enviroment

type MBR struct {
	Mbr_tamano         [10]byte
	Mbr_fecha_creacion [19]byte
	Mbr_disk_signature [10]byte
	Mbr_disk_fit       [1]byte
	Mbr_partitions     [4]Particion
}

type Particion struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [1]byte
	Part_start  [10]byte
	Part_size   [10]byte
	Part_name   [16]byte
	//Part_correlative [16]byte
	//Part_id          [4]byte
}

type EBR struct {
	Part_status [1]byte // part_mount
	Part_fit    [1]byte
	Part_start  [10]byte
	Part_size   [10]byte
	Part_next   [10]byte
	Part_name   [16]byte
}

type MountedPartitions struct {
	Path   string
	IdDisk string
	Name   string
	Index  int
	Typee  int
	Start  int
}

type UserLogged struct {
	User   string
	IdDisk string
}

type SuperBlock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [19]byte
	S_umtime            [19]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_s           int32
	S_block_s           int32
	S_first_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [16]int32
	I_type  [1]byte
	I_perm  int32
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type FolderBlock struct {
	B_content [4]Content
}

type FileBlock struct {
	B_content [64]byte
}

type ContentGraph struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	Graph string `json:"graph"`
}
