import 'package:get/get.dart';
import 'package:rupu/domain/entities/horizontal_property_voting_groups.dart';
import 'package:rupu/domain/infrastructure/repositories/horizontal_properties_repository_impl.dart';
import 'package:rupu/domain/repositories/horizontal_properties_repository.dart';

import '../horizontal_property_detail_controller.dart';

class HorizontalPropertyVotingController extends GetxController {
  final int propertyId;
  final HorizontalPropertiesRepository repository;

  HorizontalPropertyVotingController({
    required this.propertyId,
    HorizontalPropertiesRepository? repository,
  }) : repository = repository ?? HorizontalPropertiesRepositoryImpl();

  static String tagFor(int id) =>
      '${HorizontalPropertyDetailController.tagFor(id)}-voting';

  final groups = <HorizontalPropertyVotingGroup>[].obs;
  final isLoading = false.obs;
  final errorMessage = RxnString();
  final isMutationBusy = false.obs;

  int? get firstVotingGroupId =>
      groups.isNotEmpty ? groups.first.id : null;

  @override
  void onReady() {
    super.onReady();
    refresh();
  }

  Future<void> refresh() async {
    isLoading.value = true;
    errorMessage.value = null;
    try {
      final result = await repository.getHorizontalPropertyVotingGroups(
        id: propertyId,
      );
      groups.assignAll(result.groups);
      if (!result.success) {
        errorMessage.value = result.message ??
            'No se pudieron cargar los grupos de votación.';
      }
    } catch (_) {
      groups.clear();
      errorMessage.value =
          'No se pudo cargar la información de votaciones de la propiedad.';
    } finally {
      isLoading.value = false;
    }
  }

  Future<bool> createVotingGroup({
    required Map<String, dynamic> data,
  }) async {
    isMutationBusy.value = true;
    try {
      final result = await repository.createHorizontalPropertyVotingGroup(
        propertyId: propertyId,
        data: data,
      );
      if (result) {
        await refresh();
      }
      return result;
    } catch (_) {
      return false;
    } finally {
      isMutationBusy.value = false;
    }
  }

  Future<bool> updateVotingGroup({
    required int groupId,
    required Map<String, dynamic> data,
  }) async {
    isMutationBusy.value = true;
    try {
      final result = await repository.updateHorizontalPropertyVotingGroup(
        propertyId: propertyId,
        groupId: groupId,
        data: data,
      );
      if (result) {
        await refresh();
      }
      return result;
    } catch (_) {
      return false;
    } finally {
      isMutationBusy.value = false;
    }
  }

  Future<HorizontalPropertyActionResult> deleteVotingGroup(int groupId) async {
    isMutationBusy.value = true;
    try {
      final result = await repository.deleteHorizontalPropertyVotingGroup(
        propertyId: propertyId,
        groupId: groupId,
      );
      if (result.success) {
        await refresh();
      }
      return result;
    } catch (_) {
      return const HorizontalPropertyActionResult(
        success: false,
        message: 'No se pudo eliminar el grupo de votación.',
      );
    } finally {
      isMutationBusy.value = false;
    }
  }
}
